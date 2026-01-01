package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadRequest struct {
	Input struct {
		Name   string `json:"name"`
		Base64 string `json:"base64"`
	} `json:"input"`
}

type UploadResponse struct {
	URL string `json:"url"`
}

func UploadHandler(c *gin.Context) {
	var req UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			fmt.Println("❌ Error creating directory:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal directory error"})
			return
		}
	}

	
	dec, err := base64.StdEncoding.DecodeString(req.Input.Base64)
	if err != nil {
		fmt.Println("❌ Base64 Decode Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid base64 data"})
		return
	}

	
	extension := filepath.Ext(req.Input.Name)
	if extension == "" {
		extension = ".png"
	}
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	uploadPath := filepath.Join(uploadDir, newFileName)

	
	err = os.WriteFile(uploadPath, dec, 0644)
	if err != nil {
		fmt.Println("❌ File Save Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save file to disk"})
		return
	}

	
	fileURL := fmt.Sprintf("http://localhost:8082/uploads/%s", newFileName)

	fmt.Println("✅ File Uploaded Successfully:", fileURL)
	c.JSON(http.StatusOK, UploadResponse{URL: fileURL})
}