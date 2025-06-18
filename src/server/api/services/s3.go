package services

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type S3Service struct {
	client      *s3.Client
	bucketName  string
	logger      logging.Logger
	envSettings *schemas.EnvSettings
}

// Creates a new S3 service.
func NewS3Service(logger logging.Logger, envSettings *schemas.EnvSettings) *S3Service {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(envSettings.AwsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			envSettings.AwsAccessKeyId,
			envSettings.AwsSecretAccessKey,
			envSettings.AwsSessionToken,
		)),
	)
	if err != nil {
		logger.Panicln("Error loading AWS config", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:      client,
		logger:      logger,
		envSettings: envSettings,
		bucketName:  envSettings.S3BucketName,
	}
}

// Uploads a file to S3, given a key, object name and the file content.
func (s *S3Service) UploadFile(
	prefix schemas.S3Prefix,
	objectName string,
	imageBytes []byte,
) error {
	if !s.IsAvailablePrefix(prefix) {
		return fmt.Errorf("prefix %s is not available", prefix)
	}

	// Open file
	src := bytes.NewReader(imageBytes)

	// Upload the file to S3
	ctx := context.Background()
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(string(prefix) + "/" + objectName),
		Body:   src,
	})
	if err != nil {
		return err
	}

	s.logger.Infof("File successfully uploaded to S3: %s/%s", prefix, objectName)

	return nil
}

// Downloads a file from S3. It returns the file content as bytes.
func (s *S3Service) DownloadFile(prefix schemas.S3Prefix, objectName string) ([]byte, error) {
	if !s.IsAvailablePrefix(prefix) {
		return nil, fmt.Errorf("prefix %s is not available", prefix)
	}

	// Create a new context
	ctx := context.Background()

	// Download the file from S3
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(string(prefix) + "/" + objectName),
	})
	if err != nil {
		return nil, err
	}

	// Read the file content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.logger.Infof("File successfully downloaded from S3: %s/%s", prefix, objectName)

	return body, nil
}

// Checks if a prefix is available in S3.
func (s *S3Service) IsAvailablePrefix(prefix schemas.S3Prefix) bool {
	for _, p := range []schemas.S3Prefix{
		schemas.LandingPrefix,
		schemas.ProfessionalPrefix,
		schemas.LocalPrefix,
		schemas.UserPrefix,
		schemas.ServicePrefix,
		schemas.CommunityPrefix,
		schemas.TemplatePrefix,
	} {
		if p == prefix {
			return true
		}
	}

	return false
}
