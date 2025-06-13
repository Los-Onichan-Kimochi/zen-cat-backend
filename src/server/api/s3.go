package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/services"
)

// GenerateUploadPresignedURL	godoc
// @Summary 					Generate pre-signed URL for uploading images to S3
// @Description 				Generate a pre-signed URL that allows uploading images to S3 bucket
// @Tags 						S3
// @Accept 						json
// @Produce 					json
// @Param 						uploadRequest body S3UploadRequest true "Upload request details"
// @Success 					200 {object} services.PreSignedURLResponse
// @Failure 					400 {object} string "Bad request"
// @Failure 					401 {object} string "Unauthorized"
// @Failure 					403 {object} string "Forbidden"
// @Failure 					500 {object} string "Internal server error"
// @Router 						/s3/upload-url [post]
// @Security					BearerAuth
func (a *Api) GenerateUploadPresignedURL(c echo.Context) error {
	// Get user from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userRole := claims["rol"].(string)

	var request schemas.S3UploadRequest
	if err := c.Bind(&request); err != nil {
		a.Logger.Errorln("Failed to bind request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// Validate request
	if request.Folder == "" || request.FileName == "" {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "folder and fileName are required"},
		)
	}

	// Check if user has permission to upload to this folder
	// Only admins can upload images
	if userRole != string(model.UserRolAdmin) {
		return c.JSON(
			http.StatusForbidden,
			map[string]string{"error": "Only administrators can upload images"},
		)
	}

	// Initialize S3 service
	s3Service, err := services.NewS3Service(a.EnvSettings, a.Logger)
	if err != nil {
		a.Logger.Errorln("Failed to initialize S3 service:", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to connect to S3"},
		)
	}

	// Validate folder
	if !s3Service.ValidateFolder(request.Folder, true) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid folder name"})
	}

	// Add timestamp to filename to avoid conflicts
	timestampedFileName := generateTimestampedFileName(request.FileName)

	// Generate pre-signed URL
	presignedURL, err := s3Service.GenerateUploadPresignedURL(
		c.Request().Context(),
		request.Folder,
		timestampedFileName,
	)
	if err != nil {
		a.Logger.Errorln("Failed to generate upload pre-signed URL:", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to generate upload URL"},
		)
	}

	return c.JSON(http.StatusOK, presignedURL)
}

// GenerateDownloadPresignedURL	godoc
// @Summary 					Generate pre-signed URL for downloading images from S3
// @Description 				Generate a pre-signed URL that allows downloading images from S3 bucket
// @Tags 						S3
// @Accept 						json
// @Produce 					json
// @Param 						downloadRequest body S3DownloadRequest true "Download request details"
// @Success 					200 {object} services.PreSignedURLResponse
// @Failure 					400 {object} string "Bad request"
// @Failure 					401 {object} string "Unauthorized"
// @Failure 					500 {object} string "Internal server error"
// @Router 						/s3/download-url [post]
// @Security					BearerAuth
func (a *Api) GenerateDownloadPresignedURL(c echo.Context) error {
	var request schemas.S3DownloadRequest
	if err := c.Bind(&request); err != nil {
		a.Logger.Errorln("Failed to bind request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// Validate request
	if request.Key == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "key is required"})
	}

	// Initialize S3 service
	s3Service, err := services.NewS3Service(a.EnvSettings, a.Logger)
	if err != nil {
		a.Logger.Errorln("Failed to initialize S3 service:", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to connect to S3"},
		)
	}

	// Generate pre-signed URL
	presignedURL, err := s3Service.GenerateDownloadPresignedURL(c.Request().Context(), request.Key)
	if err != nil {
		a.Logger.Errorln("Failed to generate download pre-signed URL:", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to generate download URL"},
		)
	}

	return c.JSON(http.StatusOK, presignedURL)
}

// GetImageURL					godoc
// @Summary 					Get permanent image URL
// @Description 				Get the permanent URL for an image stored in S3
// @Tags 						S3
// @Accept 						json
// @Produce 					json
// @Param 						key query string true "S3 key of the image"
// @Success 					200 {object} map[string]string
// @Failure 					400 {object} string "Bad request"
// @Failure 					401 {object} string "Unauthorized"
// @Failure 					500 {object} string "Internal server error"
// @Router 						/s3/image-url [get]
// @Security					BearerAuth
func (a *Api) GetImageURL(c echo.Context) error {
	key := c.QueryParam("key")
	if key == "" {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "key parameter is required"},
		)
	}

	// Initialize S3 service
	s3Service, err := services.NewS3Service(a.EnvSettings, a.Logger)
	if err != nil {
		a.Logger.Errorln("Failed to initialize S3 service:", err)
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Failed to connect to S3"},
		)
	}

	// Get permanent image URL
	imageURL := s3Service.GetImageURL(key)

	return c.JSON(http.StatusOK, map[string]string{
		"url": imageURL,
		"key": key,
	})
}

// Helper function to generate timestamped filename
func generateTimestampedFileName(originalFileName string) string {
	timestamp := time.Now().Unix()
	id := uuid.New().String()[:8] // Use first 8 characters of UUID

	// Extract file extension
	var extension string
	for i := len(originalFileName) - 1; i >= 0; i-- {
		if originalFileName[i] == '.' {
			extension = originalFileName[i:]
			break
		}
	}

	// Remove extension from original name
	nameWithoutExt := originalFileName
	if extension != "" {
		nameWithoutExt = originalFileName[:len(originalFileName)-len(extension)]
	}

	return fmt.Sprintf("%s_%d_%s%s", nameWithoutExt, timestamp, id, extension)
}
