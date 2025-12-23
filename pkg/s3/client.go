package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

// Client wraps the S3 client with configuration
type Client struct {
	s3Client   *s3.Client
	bucketName string
	logger     *zap.Logger
}

// ClientConfig holds S3 client configuration
type ClientConfig struct {
	Region     string
	BucketName string
}

// NewClient creates a new S3 client
func NewClient(ctx context.Context, cfg ClientConfig, logger *zap.Logger) (*Client, error) {
	// Load AWS config from environment (uses AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	c := &Client{
		s3Client:   s3Client,
		bucketName: cfg.BucketName,
		logger:     logger,
	}

	logger.Info("S3 client initialized",
		zap.String("region", cfg.Region),
		zap.String("bucket", cfg.BucketName))

	return c, nil
}

// GetS3Client returns the underlying S3 client
func (c *Client) GetS3Client() *s3.Client {
	return c.s3Client
}

// GetBucketName returns the configured bucket name
func (c *Client) GetBucketName() string {
	return c.bucketName
}
