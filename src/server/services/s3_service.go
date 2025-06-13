package services

import (
	"context"
	"fmt"
	"time"

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

type PreSignedURLResponse struct {
	URL       string `json:"url"`
	Key       string `json:"key"`
	Folder    string `json:"folder"`
	ExpiresIn int64  `json:"expires_in"`
}

// NewS3Service creates a new S3 service instance
func NewS3Service(envSettings *schemas.EnvSettings, logger logging.Logger) (*S3Service, error) {
	// Create AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(envSettings.AwsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			envSettings.AwsAccessKeyId,
			envSettings.AwsSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:      client,
		bucketName:  envSettings.S3BucketName,
		logger:      logger,
		envSettings: envSettings,
	}, nil
}

// GenerateUploadPresignedURL generates a pre-signed URL for uploading files to S3
func (s *S3Service) GenerateUploadPresignedURL(
	ctx context.Context,
	folder, fileName string,
) (*PreSignedURLResponse, error) {
	// Construct the S3 key (path)
	key := fmt.Sprintf("%s/%s", folder, fileName)

	// Create presigner
	presigner := s3.NewPresignClient(s.client)

	// Set expiration time (15 minutes)
	expirationTime := 15 * time.Minute

	// Generate presigned PUT URL
	presignedRequest, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String("image/*"), // Set content type to image
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expirationTime
	})
	if err != nil {
		s.logger.Errorln("Failed to generate presigned URL for upload:", err)
		return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	s.logger.Infoln("Generated upload presigned URL for key:", key)

	return &PreSignedURLResponse{
		URL:       presignedRequest.URL,
		Key:       key,
		Folder:    folder,
		ExpiresIn: int64(expirationTime.Seconds()),
	}, nil
}

// GenerateDownloadPresignedURL generates a pre-signed URL for downloading files from S3
func (s *S3Service) GenerateDownloadPresignedURL(
	ctx context.Context,
	key string,
) (*PreSignedURLResponse, error) {
	// Create presigner
	presigner := s3.NewPresignClient(s.client)

	// Set expiration time (1 hour for downloads)
	expirationTime := 1 * time.Hour

	// Generate presigned GET URL
	presignedRequest, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expirationTime
	})
	if err != nil {
		s.logger.Errorln("Failed to generate presigned URL for download:", err)
		return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	s.logger.Infoln("Generated download presigned URL for key:", key)

	return &PreSignedURLResponse{
		URL:       presignedRequest.URL,
		Key:       key,
		ExpiresIn: int64(expirationTime.Seconds()),
	}, nil
}

// GetImageURL returns the public URL for an image (without presigning)
// This is useful when you want to store the permanent URL in the database
func (s *S3Service) GetImageURL(key string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		s.bucketName,
		s.envSettings.AwsRegion,
		key,
	)
}

// ValidateFolder checks if the folder is allowed for the current user role
func (s *S3Service) ValidateFolder(folder string, isAdmin bool) bool {
	allowedFolders := []string{
		"community",
		"landing",
		"local",
		"professional",
		"service",
		"template",
		"user",
	}

	// Check if folder is in allowed list
	for _, allowedFolder := range allowedFolders {
		if folder == allowedFolder {
			return true
		}
	}

	return false
}
