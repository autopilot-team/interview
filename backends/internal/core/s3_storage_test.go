package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testBucket     = "test-bucket"
	testAccessKey  = "minioadmin"
	testSecretKey  = "minioadmin"
	testEndpoint   = "http://localhost:9000"
	testRegion     = "us-east-1"
	testContentStr = "test content"
)

type s3TestConfig struct {
	accessKey  string
	secretKey  string
	endpoint   string
	region     string
	bucketName string
}

func getTestConfig() s3TestConfig {
	return s3TestConfig{
		accessKey:  getEnvOrDefault("AWS_ACCESS_KEY_ID", testAccessKey),
		secretKey:  getEnvOrDefault("AWS_SECRET_ACCESS_KEY", testSecretKey),
		endpoint:   getEnvOrDefault("AWS_ENDPOINT", testEndpoint),
		region:     getEnvOrDefault("AWS_REGION", testRegion),
		bucketName: testBucket,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func setupTestS3(t *testing.T) (*S3Storage, func()) {
	t.Helper()
	cfg := getTestConfig()

	// Configure S3 client
	awsCfg := aws.Config{
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.accessKey, cfg.secretKey, ""),
		Region:       cfg.region,
		BaseEndpoint: &cfg.endpoint,
	}

	// Create S3 client with ForcePathStyle for MinIO compatibility
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// Create test bucket
	ctx := context.Background()
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(cfg.bucketName),
	})
	if err != nil {
		t.Logf("Failed to create bucket, might already exist: %v", err)
	}

	storage := NewS3Storage(client, cfg.bucketName)

	cleanup := func() {
		cleanupTestBucket(t, client, cfg.bucketName)
	}

	return storage, cleanup
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

func TestS3Storage_Upload(t *testing.T) {
	storage, cleanup := setupTestS3(t)
	defer cleanup()

	tests := []struct {
		name        string
		key         string
		content     []byte
		contentType string
		metadata    map[string]string
		wantErr     bool
	}{
		{
			name:        "basic upload",
			key:         "test.txt",
			content:     []byte(testContentStr),
			contentType: "text/plain",
			metadata:    map[string]string{"test": "value"},
			wantErr:     false,
		},
		{
			name:        "empty content",
			key:         "empty.txt",
			content:     []byte{},
			contentType: "text/plain",
			wantErr:     false,
		},
		{
			name:        "with special characters in key",
			key:         "special/chars/测试.txt",
			content:     []byte(testContentStr),
			contentType: "text/plain",
			wantErr:     false,
		},
		{
			name:        "large content",
			key:         "large.txt",
			content:     bytes.Repeat([]byte("a"), 1024*1024), // 1MB
			contentType: "text/plain",
			wantErr:     false,
		},
		{
			name:        "with empty key",
			key:         "",
			content:     []byte(testContentStr),
			contentType: "text/plain",
			wantErr:     true,
		},
		{
			name:        "with nil metadata",
			key:         "nil-metadata.txt",
			content:     []byte(testContentStr),
			contentType: "text/plain",
			metadata:    nil,
			wantErr:     false,
		},
		{
			name:        "with many metadata entries",
			key:         "many-metadata.txt",
			content:     []byte(testContentStr),
			contentType: "text/plain",
			metadata: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			metadata := &ObjectMetadata{
				ContentType:    tt.contentType,
				CustomMetadata: tt.metadata,
			}

			uploadedMeta, err := storage.Upload(ctx, tt.key, bytes.NewReader(tt.content), metadata)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.key, uploadedMeta.Key)
			assert.Equal(t, tt.contentType, uploadedMeta.ContentType)
			assert.Equal(t, tt.metadata, uploadedMeta.CustomMetadata)

			// Verify content
			reader, meta, err := storage.Download(ctx, tt.key)
			require.NoError(t, err)
			defer reader.Close()

			content, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.Equal(t, tt.content, content)
			assert.Equal(t, tt.contentType, meta.ContentType)
			assert.Equal(t, tt.metadata, meta.CustomMetadata)
		})
	}
}

func TestS3Storage_List(t *testing.T) {
	storage, cleanup := setupTestS3(t)
	defer cleanup()

	tests := []struct {
		name       string
		prefix     string
		setupFiles []string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "list all",
			prefix:     "",
			setupFiles: []string{"a.txt", "b.txt", "c.txt"},
			wantCount:  3,
			wantErr:    false,
		},
		{
			name:       "list with prefix",
			prefix:     "test/",
			setupFiles: []string{"test/file1.txt", "test/file2.txt", "other.txt"},
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "no matching files",
			prefix:     "nonexistent/",
			setupFiles: []string{"test.txt"},
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "nested directories",
			prefix:     "nested/",
			setupFiles: []string{"nested/1.txt", "nested/2.txt", "nested/deep/3.txt", "nested/deep/4.txt"},
			wantCount:  4,
			wantErr:    false,
		},
		{
			name:       "with special characters in prefix",
			prefix:     "测试/",
			setupFiles: []string{"测试/file1.txt", "测试/file2.txt"},
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:   "many files",
			prefix: "many/",
			setupFiles: func() []string {
				files := make([]string, 100)
				for i := 0; i < 100; i++ {
					files[i] = fmt.Sprintf("many/file%d.txt", i)
				}
				return files
			}(),
			wantCount: 100,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Clean up any existing objects before each test case
			objects, err := storage.List(ctx, "")
			require.NoError(t, err)
			for _, obj := range objects {
				err := storage.Delete(ctx, obj.Key)
				require.NoError(t, err)
			}

			// Setup test files
			for _, key := range tt.setupFiles {
				_, err := storage.Upload(ctx, key, bytes.NewReader([]byte(testContentStr)), &ObjectMetadata{
					ContentType: "text/plain",
				})
				require.NoError(t, err)
			}

			// Test List
			objects, err = storage.List(ctx, tt.prefix)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, objects, tt.wantCount)

			// Verify all expected files are present
			foundKeys := make(map[string]bool)
			for _, obj := range objects {
				foundKeys[obj.Key] = true
			}

			for _, key := range tt.setupFiles {
				if strings.HasPrefix(key, tt.prefix) {
					assert.True(t, foundKeys[key], fmt.Sprintf("Key %s not found in list results", key))
				}
			}
		})
	}
}

func TestS3Storage_GenerateURLs(t *testing.T) {
	storage, cleanup := setupTestS3(t)
	defer cleanup()

	tests := []struct {
		name        string
		key         string
		contentType string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "generate upload URL",
			key:         "test.txt",
			contentType: "text/plain",
			expiresIn:   time.Minute,
			wantErr:     false,
		},
		{
			name:        "generate with long expiration",
			key:         "long.txt",
			contentType: "text/plain",
			expiresIn:   24 * time.Hour,
			wantErr:     false,
		},
		{
			name:        "with special characters in key",
			key:         "special/chars/测试.txt",
			contentType: "text/plain",
			expiresIn:   time.Minute,
			wantErr:     false,
		},
		{
			name:        "with empty key",
			key:         "",
			contentType: "text/plain",
			expiresIn:   time.Minute,
			wantErr:     true,
		},
		{
			name:        "with very short expiration",
			key:         "short.txt",
			contentType: "text/plain",
			expiresIn:   time.Second,
			wantErr:     false,
		},
		{
			name:        "with empty content type",
			key:         "empty-type.txt",
			contentType: "",
			expiresIn:   time.Minute,
			wantErr:     false,
		},
		{
			name:        "with binary content type",
			key:         "binary.bin",
			contentType: "application/octet-stream",
			expiresIn:   time.Minute,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Test upload URL generation
			uploadInfo, err := storage.GenerateUploadURL(ctx, tt.key, tt.contentType, tt.expiresIn)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, uploadInfo.URL)
			assert.Equal(t, "PUT", uploadInfo.Method)
			assert.Equal(t, tt.contentType, uploadInfo.Headers["Content-Type"])
			assert.True(t, uploadInfo.ExpiresAt.After(time.Now()))
			assert.True(t, uploadInfo.ExpiresAt.Before(time.Now().Add(tt.expiresIn+time.Second)))

			// Test download URL generation
			downloadInfo, err := storage.GenerateDownloadURL(ctx, tt.key, tt.expiresIn)
			require.NoError(t, err)
			assert.NotEmpty(t, downloadInfo.URL)
			assert.True(t, downloadInfo.ExpiresAt.After(time.Now()))
			assert.True(t, downloadInfo.ExpiresAt.Before(time.Now().Add(tt.expiresIn+time.Second)))
		})
	}
}

func TestS3Storage_UpdateMetadata(t *testing.T) {
	storage, cleanup := setupTestS3(t)
	defer cleanup()

	tests := []struct {
		name            string
		key             string
		initialMetadata *ObjectMetadata
		updateMetadata  *ObjectMetadata
		wantErr         bool
	}{
		{
			name: "update content type and metadata",
			key:  "test.txt",
			initialMetadata: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: map[string]string{"initial": "value"},
			},
			updateMetadata: &ObjectMetadata{
				ContentType:    "application/text",
				CustomMetadata: map[string]string{"updated": "newvalue"},
			},
			wantErr: false,
		},
		{
			name: "clear metadata",
			key:  "clear.txt",
			initialMetadata: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: map[string]string{"test": "value"},
			},
			updateMetadata: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: make(map[string]string),
			},
			wantErr: false,
		},
		{
			name: "update_non-existent_object",
			key:  "nonexistent.txt",
			updateMetadata: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: map[string]string{"test": "value"},
			},
			wantErr: true,
		},
		{
			name: "update with nil metadata",
			key:  "nil-metadata.txt",
			initialMetadata: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: map[string]string{"test": "value"},
			},
			updateMetadata: nil,
			wantErr:        true,
		},
		{
			name: "update with many metadata entries",
			key:  "many-metadata.txt",
			initialMetadata: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: map[string]string{"initial": "value"},
			},
			updateMetadata: &ObjectMetadata{
				ContentType: "text/plain",
				CustomMetadata: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
					"key4": "value4",
					"key5": "value5",
				},
			},
			wantErr: false,
		},
		{
			name: "update with special characters in metadata",
			key:  "special-metadata.txt",
			initialMetadata: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: map[string]string{"test": "value"},
			},
			updateMetadata: &ObjectMetadata{
				ContentType: "text/plain",
				CustomMetadata: map[string]string{
					"special": "测试",
					"spaces":  "value with spaces",
					"symbols": "!@#$%^&*()",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Only upload initial file if initialMetadata is provided
			if tt.initialMetadata != nil {
				_, err := storage.Upload(ctx, tt.key, bytes.NewReader([]byte(testContentStr)), tt.initialMetadata)
				require.NoError(t, err)
			}

			// Test UpdateMetadata
			err := storage.UpdateMetadata(ctx, tt.key, tt.updateMetadata)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify updated metadata
			meta, err := storage.GetMetadata(ctx, tt.key)
			require.NoError(t, err)
			assert.Equal(t, tt.updateMetadata.ContentType, meta.ContentType)
			assert.Equal(t, tt.updateMetadata.CustomMetadata, meta.CustomMetadata)
		})
	}
}

func TestS3Storage_Delete(t *testing.T) {
	storage, cleanup := setupTestS3(t)
	defer cleanup()

	tests := []struct {
		name    string
		key     string
		setup   bool
		wantErr bool
	}{
		{
			name:    "delete existing file",
			key:     "exists.txt",
			setup:   true,
			wantErr: false,
		},
		{
			name:    "delete non-existent file",
			key:     "nonexistent.txt",
			setup:   false,
			wantErr: false, // S3 delete is idempotent
		},
		{
			name:    "delete with empty key",
			key:     "",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "delete file with special characters",
			key:     "special/chars/测试.txt",
			setup:   true,
			wantErr: false,
		},
		{
			name:    "delete nested file",
			key:     "very/deeply/nested/file/structure/test.txt",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.setup {
				_, err := storage.Upload(ctx, tt.key, bytes.NewReader([]byte(testContentStr)), &ObjectMetadata{
					ContentType: "text/plain",
				})
				require.NoError(t, err)
			}

			err := storage.Delete(ctx, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file is deleted
			_, err = storage.GetMetadata(ctx, tt.key)
			assert.Error(t, err)
		})
	}
}

// Add new test for GetMetadata
func TestS3Storage_GetMetadata(t *testing.T) {
	storage, cleanup := setupTestS3(t)
	defer cleanup()

	tests := []struct {
		name        string
		key         string
		setup       bool
		setupMeta   *ObjectMetadata
		wantErr     bool
		checkFields bool
	}{
		{
			name:    "get metadata of non-existent file",
			key:     "nonexistent.txt",
			setup:   false,
			wantErr: true,
		},
		{
			name:  "get metadata of existing file",
			key:   "exists.txt",
			setup: true,
			setupMeta: &ObjectMetadata{
				ContentType:    "text/plain",
				CustomMetadata: map[string]string{"test": "value"},
			},
			wantErr:     false,
			checkFields: true,
		},
		{
			name:  "get metadata with special characters",
			key:   "special/chars/测试.txt",
			setup: true,
			setupMeta: &ObjectMetadata{
				ContentType: "text/plain",
				CustomMetadata: map[string]string{
					"special": "测试",
					"spaces":  "value with spaces",
				},
			},
			wantErr:     false,
			checkFields: true,
		},
		{
			name:    "get metadata with empty key",
			key:     "",
			setup:   false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.setup {
				_, err := storage.Upload(ctx, tt.key, bytes.NewReader([]byte(testContentStr)), tt.setupMeta)
				require.NoError(t, err)
			}

			meta, err := storage.GetMetadata(ctx, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, meta)

			if tt.checkFields {
				assert.Equal(t, tt.key, meta.Key)
				assert.Equal(t, tt.setupMeta.ContentType, meta.ContentType)
				assert.Equal(t, tt.setupMeta.CustomMetadata, meta.CustomMetadata)
				assert.NotZero(t, meta.LastModified)
				assert.NotZero(t, meta.Size)
			}
		})
	}
}
