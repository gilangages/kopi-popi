package media

import (
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadFile_InvalidExtension(t *testing.T) {
	service := NewService()

	// Simulate a PDF file upload
	file := &multipart.FileHeader{
		Filename: "document.pdf",
		Size:     1024,
	}

	url, err := service.UploadFile(file, "misc", "http://localhost:8080")

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Equal(t, "invalid file type, only JPG, JPEG, PNG, and WEBP are allowed", err.Error())
}

func TestUploadFile_SizeLimitExceeded(t *testing.T) {
	service := NewService()

	// Simulate a large image file (6MB)
	file := &multipart.FileHeader{
		Filename: "large_image.jpg",
		Size:     6 * 1024 * 1024,
	}

	url, err := service.UploadFile(file, "misc", "http://localhost:8080")

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Equal(t, "file size exceeds the 5MB limit", err.Error())
}


