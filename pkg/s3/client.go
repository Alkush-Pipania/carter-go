package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client   *s3.Client
	bucketName string
}

type ClientConfig struct {
	Region     string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
}

func NewClient(ctx context.Context, cfg ClientConfig) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKey,
				cfg.SecretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	return &Client{
		s3Client:   s3Client,
		bucketName: cfg.BucketName,
	}, nil
}

func (c *Client) GetS3Client() *s3.Client {
	return c.s3Client
}

func (c *Client) GetBucketName() string {
	return c.bucketName
}
