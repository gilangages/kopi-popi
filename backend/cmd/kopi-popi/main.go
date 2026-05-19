package main

import (
	config "github.com/gilangages/kopi-popi/configs"
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Koneksi ke Database
	// Memastikan database menyala sebelum route dijalankan
	config.ConnectDB()

	// 2. Setup Framework Gin (Router)
	r := gin.Default()

	// 3. Setup Global Middleware CORS (Opsional tapi wajib untuk dipanggil frontend/React/Vue nanti)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 4. Register Health Check Endpoint (Public)
	// Kita tes panggil Response Helper buatan kita di Fase 1
	r.GET("/", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"message": "Welcome to Kopi-Popi API!",
			"version": "1.0.0",
		})
	})

	// Todo di Fase Selanjutnya:
	// Di sini nanti kita akan mendaftarkan router per-domain.
	// Contoh: 
	// auth_routes.Setup(r)
	// catalog_routes.Setup(r)

	// 5. Jalankan Server di port 8080
	r.Run(":8080")
}
