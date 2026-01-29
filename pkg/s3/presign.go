package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

// Presigner handles generating presigned URLs for S3 operations
type Presigner struct {
	presignClient *s3.PresignClient
	bucketName    string
	expiry        time.Duration
	logger        *zap.Logger
}

// PresignerConfig holds presigner configuration
type PresignerConfig struct {
	ExpiryMinutes int
}

// NewPresigner creates a new presigner from an S3 client
func NewPresigner(client *Client, cfg PresignerConfig, logger *zap.Logger) *Presigner {
	presignClient := s3.NewPresignClient(client.GetS3Client())

	return &Presigner{
		presignClient: presignClient,
		bucketName:    client.GetBucketName(),
		expiry:        time.Duration(cfg.ExpiryMinutes) * time.Minute,
		logger:        logger,
	}
}

// PresignedUploadResult contains the presigned URL and metadata
type PresignedUploadResult struct {
	URL       string    `json:"url"`
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GenerateUploadURL creates a presigned PUT URL for uploading a file
func (p *Presigner) GenerateUploadURL(ctx context.Context, key string, contentType string) (*PresignedUploadResult, error) {
	request, err := p.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(p.expiry))
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	expiresAt := time.Now().Add(p.expiry)

	p.logger.Debug("Generated presigned upload URL",
		zap.String("key", key),
		zap.String("content_type", contentType),
		zap.Time("expires_at", expiresAt))

	return &PresignedUploadResult{
		URL:       request.URL,
		Key:       key,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateDownloadURL creates a presigned GET URL for downloading a file
func (p *Presigner) GenerateDownloadURL(ctx context.Context, key string) (*PresignedUploadResult, error) {
	request, err := p.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = p.expiry
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	expiresAt := time.Now().Add(p.expiry)

	p.logger.Debug("Generated presigned download URL",
		zap.String("key", key),
		zap.Time("expires_at", expiresAt))

	return &PresignedUploadResult{
		URL:       request.URL,
		Key:       key,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateS3Key creates a consistent S3 key for source files
// Format: sources/{user_id}/{source_id}/{filename}
func GenerateS3Key(userID, sourceID, filename string) string {
	return fmt.Sprintf("sources/%s/%s/%s", userID, sourceID, filename)
}
