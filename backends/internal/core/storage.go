package core

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

// ObjectMetadata represents metadata for a stored object
type ObjectMetadata struct {
	// Key is the unique identifier for the object
	Key string `json:"key"`

	// Size in bytes
	Size int64 `json:"size"`

	// ContentType of the object
	ContentType string `json:"content_type"`

	// ETag for object versioning
	ETag string `json:"etag,omitempty"`

	// LastModified timestamp
	LastModified time.Time `json:"last_modified"`

	// Custom metadata as key-value pairs
	CustomMetadata map[string]string `json:"custom_metadata,omitempty"`
}

// UploadInfo contains information needed for direct uploads
type UploadInfo struct {
	// URL for the direct upload
	URL string `json:"url"`

	// Method to use for the upload (PUT/POST)
	Method string `json:"method"`

	// Headers required for the upload
	Headers map[string]string `json:"headers,omitempty"`

	// ExpiresAt indicates when the upload URL expires
	ExpiresAt time.Time `json:"expires_at"`
}

// DownloadInfo contains information needed for direct downloads
type DownloadInfo struct {
	// URL for the direct download
	URL string `json:"url"`

	// Headers required for the download
	Headers map[string]string `json:"headers,omitempty"`

	// ExpiresAt indicates when the download URL expires
	ExpiresAt time.Time `json:"expires_at"`
}

// Storage defines the interface for object storage operations
type Storage interface {
	// Upload stores an object with the given key and returns its metadata
	Upload(ctx context.Context, key string, reader io.Reader, metadata *ObjectMetadata) (*ObjectMetadata, error)

	// Download retrieves an object and its metadata by key
	Download(ctx context.Context, key string) (io.ReadCloser, *ObjectMetadata, error)

	// Delete removes an object by key
	Delete(ctx context.Context, key string) error

	// GetMetadata retrieves just the metadata for an object
	GetMetadata(ctx context.Context, key string) (*ObjectMetadata, error)

	// List returns metadata for objects with the given prefix
	List(ctx context.Context, prefix string) ([]*ObjectMetadata, error)

	// GenerateUploadURL creates a pre-signed URL for direct upload
	GenerateUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (*UploadInfo, error)

	// GenerateDownloadURL creates a pre-signed URL for direct download
	GenerateDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (*DownloadInfo, error)

	// UpdateMetadata updates the metadata for an existing object
	UpdateMetadata(ctx context.Context, key string, metadata *ObjectMetadata) error
}

type StorageOptions struct {
	// Logger is used for storage access-related logging
	Logger *slog.Logger
	// Endpoint is used for public accessible endpoints
	Endpoint string
	// Region is the location of the storge server.
	Region string
	// AccessKeyID is the storage access key.
	AccessKeyID string
	// AccessKeyID is the storage secret key.
	SecretAccessKey string
	// Bucket is the bucket name to use.
	Bucket string
	// UsePathStyle enables client path style addressing.
	UsePathStyle bool
}

func NewStorage(ctx context.Context, opts StorageOptions) (Storage, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	awsCfg := aws.Config{
		Logger:       compatLogger{opts.Logger},
		Credentials:  credentials.NewStaticCredentialsProvider(opts.AccessKeyID, opts.SecretAccessKey, ""),
		Region:       opts.Region,
		BaseEndpoint: &opts.Endpoint,
	}
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = opts.UsePathStyle // Required for MinIO compat
	})

	return NewS3Storage(s3Client, opts.Bucket), nil
}

type compatLogger struct {
	*slog.Logger
}

func (c compatLogger) Logf(typ logging.Classification, format string, v ...any) {
	switch typ {
	case logging.Warn:
		c.Warn(fmt.Sprintf(format, v...))
	case logging.Debug:
		c.Debug(fmt.Sprintf(format, v...))
	}
}
