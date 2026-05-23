package main

import (
	config "github.com/gilangages/kopi-popi/configs"
	"github.com/gilangages/kopi-popi/internal/auth"
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Koneksi ke Database
	// Memastikan database menyala sebelum route dijalankan
	db := config.ConnectDB()
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// 2. Setup Framework Gin (Router)
	r := gin.Default()

	// 3. Setup Global Middleware CORS
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
	r.GET("/", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"message": "Welcome to Kopi-Popi API!",
			"version": "1.0.0",
		})
	})

	// 5. Inisialisasi Domain Auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// 6. Daftarkan router per-domain
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		authRoutes.DELETE("/logout", authHandler.Logout)
	}

	// 7. Jalankan Server di port 8080
	r.Run(":8080")
}
