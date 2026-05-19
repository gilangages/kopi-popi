package middleware

import (
	"net/http"

	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

// RequireRole mengecek apakah role user (yang diambil dari JWT) sesuai dengan allowedRoles.
// WAJIB dipanggil setelah RequireAuth(), karena kita butuh data "role" dari context.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context (yang diset oleh RequireAuth)
		userRole, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusForbidden, "Role not found in context")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "Failed to parse role")
			c.Abort()
			return
		}

		// Cek apakah roleStr termasuk dalam list yang diizinkan
		isAllowed := false
		for _, role := range allowedRoles {
			if roleStr == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response.Error(c, http.StatusForbidden, "You don't have permission to access this endpoint")
			c.Abort()
			return
		}

		// Lolos dari middleware, lanjut ke handler berikutnya
		c.Next()
	}
}
