package middleware

import (
	"strings"

	"github.com/gilangages/kopi-popi/pkg/jwt"
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

// RequireAuth adalah middleware untuk mengecek keabsahan token JWT dari header Authorization.
// Jika valid, data payload (claims) akan disisipkan ke dalam context (c.Set).
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, "Unauthorized: Missing Authorization header")
			c.Abort()
			return
		}

		// Header harus berformat "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, "Unauthorized: Invalid Authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, 401, "Unauthorized: Invalid or expired token")
			c.Abort()
			return
		}

		// Menyisipkan data claims ke dalam Gin context untuk dipakai oleh Handler selanjutnya
		c.Set("user_id", claims["user_id"])
		c.Set("name", claims["name"])
		c.Set("role", claims["role"])
		c.Set("branch_id", claims["branch_id"]) // Bisa null (nil) untuk Admin/Customer

		c.Next()
	}
}
