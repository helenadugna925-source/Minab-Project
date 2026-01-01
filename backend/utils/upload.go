package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	
)

// Ensure the upload directory exists
func InitUploadPath() {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}
}

func SaveBase64Image(base64Str string, name string) (string, error) {
	// 1. Separate metadata from data (e.g., "data:image/png;base64,.....")
	parts := strings.Split(base64Str, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 format")
	}

	// 2. Decode
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// 3. Generate Unique Filename
	ext := filepath.Ext(name)
	if ext == "" {
		ext = ".jpg" // Default fallback
	}
	newFileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)
	filePath := filepath.Join("./uploads", newFileName)

	// 4. Save File
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", err
	}

	// Return the public URL (Assuming your server runs on port 8082)
	return fmt.Sprintf("http://localhost:8082/uploads/%s", newFileName), nil
}