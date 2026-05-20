package auth

import (
	"net/http"

	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		// Asumsi sederhana untuk status code error
		statusCode := http.StatusInternalServerError
		if err.Error() == "email is already registered" {
			statusCode = http.StatusConflict
		}
		response.Error(c, statusCode, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
