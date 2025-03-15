package core

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements the Storage interface for AWS S3
type S3Storage struct {
	client     *s3.Client
	bucketName string
}

// NewS3Storage creates a new S3Storage instance
func NewS3Storage(client *s3.Client, bucketName string) *S3Storage {
	return &S3Storage{
		client:     client,
		bucketName: bucketName,
	}
}

// Upload implements Storage.Upload
func (s *S3Storage) Upload(ctx context.Context, key string, reader io.Reader, metadata *ObjectMetadata) (*ObjectMetadata, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(metadata.ContentType),
		Metadata:    metadata.CustomMetadata,
	}

	result, err := s.client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload object: %w", err)
	}

	return &ObjectMetadata{
		Key:            key,
		ContentType:    metadata.ContentType,
		ETag:           aws.ToString(result.ETag),
		LastModified:   time.Now(),
		CustomMetadata: metadata.CustomMetadata,
	}, nil
}

// Download implements Storage.Download
func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, *ObjectMetadata, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download object: %w", err)
	}

	metadata := &ObjectMetadata{
		Key:            key,
		Size:           aws.ToInt64(result.ContentLength),
		ContentType:    aws.ToString(result.ContentType),
		ETag:           aws.ToString(result.ETag),
		LastModified:   aws.ToTime(result.LastModified),
		CustomMetadata: result.Metadata,
	}

	return result.Body, metadata, nil
}

// Delete implements Storage.Delete
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// GetMetadata implements Storage.GetMetadata
func (s *S3Storage) GetMetadata(ctx context.Context, key string) (*ObjectMetadata, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	result, err := s.client.HeadObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	// Initialize an empty map if Metadata is nil
	metadata := result.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	return &ObjectMetadata{
		Key:            key,
		Size:           aws.ToInt64(result.ContentLength),
		ContentType:    aws.ToString(result.ContentType),
		ETag:           aws.ToString(result.ETag),
		LastModified:   aws.ToTime(result.LastModified),
		CustomMetadata: metadata,
	}, nil
}

// List implements Storage.List
func (s *S3Storage) List(ctx context.Context, prefix string) ([]*ObjectMetadata, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(prefix),
	}

	var objects []*ObjectMetadata
	paginator := s3.NewListObjectsV2Paginator(s.client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			objects = append(objects, &ObjectMetadata{
				Key:          aws.ToString(obj.Key),
				Size:         aws.ToInt64(obj.Size),
				ETag:         aws.ToString(obj.ETag),
				LastModified: aws.ToTime(obj.LastModified),
			})
		}
	}

	return objects, nil
}

// GenerateUploadURL implements Storage.GenerateUploadURL
func (s *S3Storage) GenerateUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (*UploadInfo, error) {
	presignClient := s3.NewPresignClient(s.client)
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	result, err := presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return &UploadInfo{
		URL:       result.URL,
		Method:    string(result.Method),
		Headers:   map[string]string{"Content-Type": contentType},
		ExpiresAt: time.Now().Add(expiresIn),
	}, nil
}

// GenerateDownloadURL implements Storage.GenerateDownloadURL
func (s *S3Storage) GenerateDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (*DownloadInfo, error) {
	presignClient := s3.NewPresignClient(s.client)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	var opts []func(*s3.PresignOptions)
	if expiresIn == 0 {
		opts = append(opts, s3.WithPresignExpires(expiresIn))
	}
	result, err := presignClient.PresignGetObject(ctx, input, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}

	return &DownloadInfo{
		URL:       result.URL,
		Headers:   make(map[string]string),
		ExpiresAt: time.Now().Add(expiresIn),
	}, nil
}

// UpdateMetadata implements Storage.UpdateMetadata
func (s *S3Storage) UpdateMetadata(ctx context.Context, key string, metadata *ObjectMetadata) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	// Ensure CustomMetadata is initialized
	if metadata.CustomMetadata == nil {
		metadata.CustomMetadata = make(map[string]string)
	}

	// S3 requires copying the object to itself to update metadata
	copySource := fmt.Sprintf("%s/%s", s.bucketName, key)
	input := &s3.CopyObjectInput{
		Bucket:            aws.String(s.bucketName),
		CopySource:        aws.String(copySource),
		Key:               aws.String(key),
		ContentType:       aws.String(metadata.ContentType),
		Metadata:          metadata.CustomMetadata,
		MetadataDirective: types.MetadataDirectiveReplace,
	}

	_, err := s.client.CopyObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update object metadata: %w", err)
	}

	return nil
}
