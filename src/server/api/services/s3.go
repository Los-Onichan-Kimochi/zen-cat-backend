package services

import (
	"bytes"
	"context"
	"errors"
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
		s.logger.Errorf("Prefix %s is not available", prefix)
		return errors.New("prefix not available")
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
		s.logger.Errorf("Error uploading file to S3: %v", err)
		return err
	}

	s.logger.Infof("File %s uploaded to S3.", string(prefix)+"/"+objectName)

	return nil
}

// Downloads a file from S3. It returns the file content as bytes.
func (s *S3Service) DownloadFile(prefix schemas.S3Prefix, objectName string) ([]byte, error) {
	if !s.IsAvailablePrefix(prefix) {
		s.logger.Errorf("Prefix %s is not available", prefix)
		return nil, errors.New("prefix not available")
	}

	// Create a new context
	ctx := context.Background()

	// Download the file from S3
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(string(prefix) + "/" + objectName),
	})
	if err != nil {
		s.logger.Errorf("Error downloading file %s from S3: %v", string(prefix)+"/"+objectName, err)
		return nil, err
	}

	// Read the file content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Error reading file content: %v", err)
		return nil, err
	}

	s.logger.Infof("File %s downloaded from S3.", string(prefix)+"/"+objectName)

	return body, nil
}

// Checks if a prefix is available in S3.
func (s *S3Service) IsAvailablePrefix(prefix schemas.S3Prefix) bool {
	for _, p := range []schemas.S3Prefix{
		schemas.LandingS3Prefix,
		schemas.ProfessionalS3Prefix,
		schemas.LocalS3Prefix,
		schemas.UserS3Prefix,
		schemas.ServiceS3Prefix,
		schemas.CommunityS3Prefix,
		schemas.TemplateS3Prefix,
	} {
		if p == prefix {
			return true
		}
	}

	return false
}
