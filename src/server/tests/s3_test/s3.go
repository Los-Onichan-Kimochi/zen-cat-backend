package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/api/services"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

func main() {
	logger := logging.NewLogger("S3 Test", "Version 1.0", logging.FormatText, 4)
	s3Service := services.NewS3Service(
		logger,
		schemas.NewEnvSettings(logger),
	)

	// Crear directorio download si no existe
	downloadDir := "download"
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating download directory: %v\n", err)
		return
	}

	// Lista de archivos para probar
	testFiles := []string{
		"mitos_embarazo_-t.jpg",
		"gym_rats.webp",
	}

	fmt.Println("🚀 Starting S3 test...")
	fmt.Printf("📁 Download directory: %s\n", downloadDir)
	fmt.Println("=" + fmt.Sprintf("%40s", "="))

	// Probar subida y descarga para cada archivo
	for _, filePath := range testFiles {
		fmt.Printf("\n📄 Testing file: %s\n", filePath)

		// Verificar que el archivo existe
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("⚠️  File %s not found, skipping...\n", filePath)
			continue
		}

		// UploadFile
		fmt.Println("⬆️  Uploading file...")
		if err := uploadFile(s3Service, filePath); err != nil {
			fmt.Printf("❌ Upload failed: %v\n", err)
			continue
		}

		// DownloadFile
		fmt.Println("⬇️  Downloading file...")
		if err := downloadFile(s3Service, filePath, downloadDir); err != nil {
			fmt.Printf("❌ Download failed: %v\n", err)
			continue
		}

		fmt.Printf("✅ Test completed for %s\n", filePath)
	}

	fmt.Println("\n🎉 S3 test successfully completed!")
}

func downloadFile(s3Service *services.S3Service, fileName string, downloadDir string) error {
	imageBytes, err := s3Service.DownloadFile(schemas.LandingPrefix, fileName)
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}

	// Crear la ruta completa para el archivo descargado
	downloadPath := filepath.Join(downloadDir, fileName)

	err = os.WriteFile(downloadPath, imageBytes, 0o644)
	if err != nil {
		return fmt.Errorf("error writing file to %s: %w", downloadPath, err)
	}

	fmt.Printf("   ✅ File downloaded to: %s\n", downloadPath)
	return nil
}

func uploadFile(s3Service *services.S3Service, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	err = s3Service.UploadFile(schemas.LandingPrefix, filePath, imageBytes)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	fmt.Printf("   ✅ File uploaded to S3: %s\n", filePath)
	return nil
}
