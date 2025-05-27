package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	
	"github.com/mingeme/imgood/internal/config"
)

// Client represents an S3 client
type Client struct {
	s3Client *s3.Client
	config   config.S3Config
}

// NewClient creates a new S3 client with the provided configuration
func NewClient(cfg config.S3Config) (*Client, error) {
	// Configure AWS
	awsCfg, err := configureAWS(cfg.Region, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error configuring AWS: %w", err)
	}

	// Create S3 client with options
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	return &Client{
		s3Client: s3Client,
		config:   cfg,
	}, nil
}

// UploadFile uploads a file to S3
func (c *Client) UploadFile(key string, data []byte) error {
	ctx := context.Background()
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("error uploading to S3: %w", err)
	}

	return nil
}

// GetFileURL returns the URL for an uploaded file
func (c *Client) GetFileURL(key string) string {
	if c.config.Endpoint != "" {
		// For custom S3 endpoints
		return fmt.Sprintf("%s/%s/%s", c.config.Endpoint, c.config.Bucket, key)
	} 
	// For AWS S3
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.config.Bucket, c.config.Region, key)
}

// configureAWS sets up the AWS configuration with the provided credentials and region
func configureAWS(region, accessKey, secretKey string) (aws.Config, error) {
	configOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
	}

	// Add credentials if provided
	if accessKey != "" && secretKey != "" {
		configOptions = append(configOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}

	// Load the AWS configuration
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), configOptions...)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
