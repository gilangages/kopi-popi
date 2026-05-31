package media

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Service interface {
	UploadFile(file *multipart.FileHeader, folder string, baseURL string) (string, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) UploadFile(file *multipart.FileHeader, folder string, baseURL string) (string, error) {
	// 1. Validasi Ekstensi (Hanya Gambar)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return "", errors.New("invalid file type, only JPG, JPEG, PNG, and WEBP are allowed")
	}

	// 2. Validasi Ukuran File (Maksimal 5MB)
	const maxUploadSize = 5 * 1024 * 1024 // 5 MB
	if file.Size > maxUploadSize {
		return "", errors.New("file size exceeds the 5MB limit")
	}

	// 3. Sanitasi & Penamaan File Baru
	// Menggunakan timestamp agar tidak bentrok (contoh: 1689328293-image.jpg)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// 4. Pastikan direktori tujuan ada (Auto-create folder)
	uploadPath := filepath.Join("uploads", folder)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// 5. Simpan file fisik
	dstPath := filepath.Join(uploadPath, newFileName)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// 6. Kembalikan Public URL
	// Gunakan forward slash (/) untuk URL, bukan backslash (\)
	// Output URL contoh: http://localhost:8080/uploads/products/1689328293.jpg
	publicURL := fmt.Sprintf("%s/uploads/%s/%s", baseURL, folder, newFileName)
	return publicURL, nil
}
