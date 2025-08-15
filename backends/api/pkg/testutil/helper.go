package testutil

import (
	"autopilot/backends/api/internal"
	"autopilot/backends/api/internal/identity"
	"autopilot/backends/api/internal/identity/service"
	"autopilot/backends/api/internal/payment"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/app/mocks"
	"autopilot/backends/api/pkg/middleware"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorker is a mock implementation of core.Worker using testify/mock.
type MockWorker struct {
	mock.Mock

	MailRequests []service.MailerArgs
}

// GetClient implements core.Worker
func (m *MockWorker) GetClient() *river.Client[pgx.Tx] {
	// Return nil or a mock client depending on your testing needs
	return nil
}

// GetDBPool implements core.Worker
func (m *MockWorker) GetDBPool() *pgxpool.Pool {
	// Return nil or a mock client depending on your testing needs
	return nil
}

// Insert implements core.Worker
func (m *MockWorker) Insert(ctx context.Context, job river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	switch v := job.(type) {
	case service.MailerArgs:
		m.MailRequests = append(m.MailRequests, v)
	}
	args := m.Called(ctx, job, opts)
	return nil, args.Error(1)
}

// Queues implements core.Worker
func (m *MockWorker) Queues() *river.QueueBundle {
	return nil
}

// Reset cleans up internal cache for subsequent tasks.
func (m *MockWorker) Reset() {
	m.MailRequests = nil
}

// Start implements core.Worker
func (m *MockWorker) Start(ctx context.Context) error {
	return nil
}

// Stop implements core.Worker
func (m *MockWorker) Stop(ctx context.Context) error {
	return nil
}

// StopAndCancel implements core.Worker
func (m *MockWorker) StopAndCancel(ctx context.Context) error {
	return nil
}

// Stopped implements core.Worker
func (m *MockWorker) Stopped() <-chan struct{} {
	return nil
}

// Container creates a new mock container for testing
func Container(t *testing.T) (humatest.TestAPI, *app.Container, *internal.Module) {
	ctx := context.Background()
	cleanUp := make([]func() error, 0)
	mode := types.DebugMode

	// Initialize the configuration
	config, err := app.NewConfig()
	assert.NoError(t, err)

	// Initialize the logger
	logger := core.NewLogger(core.LoggerOptions{Mode: mode})

	// Initialize the local filesystem
	localFS, err := core.NewLocalFS("./backends/api")
	assert.NoError(t, err)

	// Initialize the i18n bundle
	i18nBundle, err := core.NewI18nBundle(localFS, "locales")
	assert.NoError(t, err)

	// Initialize the mailer
	mailer, err := core.NewMail(core.MailOptions{
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       mode,
		PreviewData: map[string]map[string]any{
			"welcome": {
				"AppName":         config.App.Name,
				"AssetsURL":       config.App.AssetsURL,
				"Duration":        (24 * time.Hour).Hours(),
				"Name":            "John Doe",
				"VerificationURL": fmt.Sprintf("%s/verify-email?token=01948450-988e-7976-a454-7163b6f1c6c6", config.App.DashboardURL),
			},
		},
		SMTPURL: config.App.Mailer.SMTPURL,
		TemplateOptions: &core.MailTemplateOptions{
			Dir: "templates",
			ExtraFuncs: []template.FuncMap{
				{
					"currentYear": func() string {
						return strconv.Itoa(time.Now().Year())
					},
				},
			},
			FS:     localFS,
			Layout: "layouts/transactional",
		},
	})
	assert.NoError(t, err)

	projectRoot, err := core.FindProjectRoot()
	assert.NoError(t, err)
	mainFile := filepath.Join(projectRoot, "backends/api/internal/app/main.go")

	// Initialize the identity database
	identityDB, identityDBCleanUp, err := NewDB(ctx, "identity", core.DBOptions{
		Identifier:        "identity",
		Logger:            logger,
		DisableLogQueries: true,
		MainFile:          mainFile,
		MigrationsFS:      localFS,
	})
	assert.NoError(t, err)
	cleanUp = append(cleanUp, func() error {
		identityDBCleanUp()
		return nil
	})

	// Initialize the payment databases
	paymentLiveDB, paymentLiveDBCleanUp, err := NewDB(ctx, "payment_live", core.DBOptions{
		Identifier:        "payment",
		Logger:            logger,
		DisableLogQueries: true,
		MainFile:          mainFile,
		MigrationsFS:      localFS,
	})
	assert.NoError(t, err)
	cleanUp = append(cleanUp, func() error {
		paymentLiveDBCleanUp()
		return nil
	})

	paymentTestDB, paymentTestDBCleanUp, err := NewDB(ctx, "payment_test", core.DBOptions{
		Identifier:        "payment",
		Logger:            logger,
		DisableLogQueries: true,
		MainFile:          mainFile,
		MigrationsFS:      localFS,
	})
	assert.NoError(t, err)
	cleanUp = append(cleanUp, func() error {
		paymentTestDBCleanUp()
		return nil
	})

	mockWorker := &MockWorker{}
	mockWorker.On("Insert", mock.Anything, mock.AnythingOfType("MailerArgs"), mock.Anything).
		Return(nil, nil)

	container := &app.Container{
		CleanUp: cleanUp,
		Config:  config,
		DB: app.ContainerDB{
			Identity: identityDB,
			Payment: app.PaymentDB{
				Live: paymentLiveDB,
				Test: paymentTestDB,
			},
		},
		FS: app.ContainerFS{
			Locales:    localFS,
			Migrations: localFS,
			Templates:  localFS,
		},
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       mode,
		Mailer:     mailer,
		Storage: app.ContainerS3{
			Identity: newS3(t),
		},
		Turnstile: &mocks.MockTurnstiler{},
		Worker:    mockWorker,
	}

	identityMod, err := identity.New(ctx, container)
	assert.NoError(t, err)

	paymentMod, err := payment.New(ctx, container)
	assert.NoError(t, err)

	_, api := humatest.New(t)
	api.UseMiddleware(
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithContext(ctx, middleware.AttachRequestMetadata(ctx.Context(),
				ctx.Header("X-Forwarded-For"),
				ctx.RemoteAddr(),
				ctx.Header("User-Agent"),
				ctx.Header(middleware.CFCountryHeader)))
			next(ctx)
		},

		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithContext(ctx, middleware.AttachContainer(ctx.Context(), container))
			ctx = huma.WithContext(ctx, context.WithValue(ctx.Context(), middleware.EntityKey, ctx.Header(middleware.ActiveEntityHeader)))
			next(ctx)
		},
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithContext(ctx, middleware.AttachT(ctx.Context(), container.I18nBundle, ctx.Header("Accept-Language"), ctx.Query("locale")))
			next(ctx)
		},
		func(ctx huma.Context, next func(huma.Context)) {
			mode := types.OperationModeTest
			if hMode := ctx.Header("X-Operation-Mode"); hMode != "" {
				if hMode == string(types.OperationModeLive) {
					mode = types.OperationModeLive
				}
			}
			ctx = huma.WithContext(ctx, context.WithValue(ctx.Context(), types.OperationModeKey, mode))
			next(ctx)
		},
	)

	t.Cleanup(func() {
		container.Close()
	})
	return api, container, &internal.Module{
		Identity: identityMod,
		Payment:  paymentMod,
	}
}

const (
	testBucket    = "test-bucket"
	testAccessKey = "minioadmin"
	testSecretKey = "minioadmin"
	testRegion    = "us-east-1"
)

var testEndpoint = "http://localhost:9000"

func newS3(t *testing.T) *core.S3Storage {
	t.Helper()
	// Configure S3 client
	awsCfg := aws.Config{
		Credentials:  credentials.NewStaticCredentialsProvider(testAccessKey, testSecretKey, ""),
		Region:       testRegion,
		BaseEndpoint: &testEndpoint,
	}

	// Create S3 client with ForcePathStyle for MinIO compatibility
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	bucketName := uuid.New().String()
	// Create test bucket
	ctx := context.Background()
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Logf("Failed to create bucket, might already exist: %v", err)
	}

	storage := core.NewS3Storage(client, bucketName)
	t.Cleanup(func() {
		cleanupTestBucket(t, client, bucketName)
	})
	return storage
}

func cleanupTestBucket(t *testing.T, client *s3.Client, bucketName string) {
	t.Helper()
	ctx := context.Background()

	// List and delete all objects
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
	paginator := s3.NewListObjectsV2Paginator(client, listInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			t.Logf("Failed to list objects: %v", err)
			return
		}

		for _, obj := range page.Contents {
			_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    obj.Key,
			})
			if err != nil {
				t.Logf("Failed to delete object %s: %v", *obj.Key, err)
			}
		}
	}

	// Delete bucket
	_, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Logf("Failed to delete bucket: %v", err)
	}
}
