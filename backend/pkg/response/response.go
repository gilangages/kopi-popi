package response

import "github.com/gin-gonic/gin"

// SuccessResponse adalah format standar untuk response sukses sesuai OpenAPI kita
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse adalah format standar untuk response error sesuai OpenAPI kita
type ErrorResponse struct {
	Errors ErrorDetail `json:"errors"`
}

// ErrorDetail menyimpan detail pesan error
type ErrorDetail struct {
	Message string `json:"message"`
}

// Success adalah helper untuk mengembalikan respons JSON sukses
// Contoh penggunaan: response.Success(c, http.StatusOK, gin.H{"user": user})
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Data: data,
	})
}

// Error adalah helper untuk mengembalikan respons JSON gagal/error
// Contoh penggunaan: response.Error(c, http.StatusBadRequest, "Invalid email or password")
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Errors: ErrorDetail{
			Message: message,
		},
	})
}
