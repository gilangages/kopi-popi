package media

import (
	"fmt"
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) UploadFile(c *gin.Context) {
	// Parse input file form "file"
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "file is required")
		return
	}

	// Parse input text form "folder", default ke "misc" jika kosong
	folder := c.PostForm("folder")
	if folder == "" {
		folder = "misc"
	}

	// Tentukan Base URL server
	// Dalam praktiknya bisa diambil dari Env Var (misal: APP_URL).
	// Di sini kita kombinasikan skema HTTP dan host dari request asli
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	// Lakukan proses unggah lewat Service
	publicURL, err := h.service.UploadFile(file, folder, baseURL)
	if err != nil {
		// Mengembalikan bad request jika ada masalah ekstensi atau ukuran
		response.Error(c, 400, err.Error())
		return
	}

	// Kembalikan URL sukses
	c.JSON(200, gin.H{
		"code":    200,
		"message": "file uploaded successfully",
		"data": gin.H{
			"url": publicURL,
		},
	})
}
