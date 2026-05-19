package middleware

import (
	"net/http"
	"strings"

	"github.com/gilangages/kopi-popi/pkg/jwt"
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

// RequireAuth adalah penjaga gerbang yang mewajibkan user mengirim JWT Token.
// Dipasang di endpoint yang bersifat "Protected".
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is missing")
			c.Abort()
			return
		}

		// Format token harus: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format. Use 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi signature token menggunakan package jwt kita
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Simpan data (claims) dari token ke dalam Context
		// Tujuannya agar nanti di Handler kita bisa memanggil c.Get("user_id")
		c.Set("user_id", claims["user_id"])
		c.Set("name", claims["name"])
		c.Set("role", claims["role"])

		// Lolos dari middleware, lanjut ke handler berikutnya
		c.Next()
	}
}
