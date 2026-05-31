package user

import (
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// GetMyProfile
func (h *Handler) GetMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	profile, err := h.service.GetMyProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, profile)
}

// UpdateMyProfile
func (h *Handler) UpdateMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	if err := h.service.UpdateMyProfile(c.Request.Context(), userID, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "profile updated successfully")
}

// DeleteProfilePicture
func (h *Handler) DeleteProfilePicture(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.service.DeleteProfilePicture(c.Request.Context(), userID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "profile picture deleted successfully")
}

// UpdateMyPassword
func (h *Handler) UpdateMyPassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	if err := h.service.UpdateMyPassword(c.Request.Context(), userID, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "password updated successfully")
}

// RequestEmailOTP
func (h *Handler) RequestEmailOTP(c *gin.Context) {
	userID := c.GetString("user_id")

	var req RequestEmailOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	if err := h.service.RequestEmailOTP(c.Request.Context(), userID, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "OTP has been sent to new email")
}

// VerifyEmailOTP
func (h *Handler) VerifyEmailOTP(c *gin.Context) {
	userID := c.GetString("user_id")

	var req VerifyEmailOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	if err := h.service.VerifyEmailOTP(c.Request.Context(), userID, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "email updated successfully")
}

// CreateManager (Admin Only)
func (h *Handler) CreateManager(c *gin.Context) {
	// Middleware Otorisasi Manual (Bisa juga dipisah di middleware khusus Role)
	if c.GetString("role") != "Admin" {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}

	var req CreateManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	if err := h.service.CreateManager(c.Request.Context(), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "manager created successfully")
}

// CreateCashier (Manager Only)
func (h *Handler) CreateCashier(c *gin.Context) {
	if c.GetString("role") != "Manager" {
		response.Error(c, 403, "forbidden: manager role required")
		return
	}

	// Ambil BranchID dari JWT (Float64 karena JWT MapClaims secara default unmarshal angka ke float64)
	val, ok := c.Get("branch_id")
	if !ok || val == nil {
		response.Error(c, 403, "forbidden: manager must be assigned to a branch")
		return
	}
	managerBranchID := int(val.(float64))

	var req CreateCashierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	if err := h.service.CreateCashier(c.Request.Context(), managerBranchID, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "cashier created successfully")
}

// GetEmployees
func (h *Handler) GetEmployees(c *gin.Context) {
	role := c.GetString("role")

	// Cek Role
	if role != "Admin" && role != "Manager" {
		response.Error(c, 403, "forbidden: insufficient privileges")
		return
	}

	var branchIDPtr *int
	val, ok := c.Get("branch_id")
	if ok && val != nil {
		bID := int(val.(float64))
		branchIDPtr = &bID
	}

	employees, err := h.service.GetEmployees(c.Request.Context(), role, branchIDPtr)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, employees)
}

// DisableEmployee (Admin Only)
func (h *Handler) DisableEmployee(c *gin.Context) {
	role := c.GetString("role")
	if role != "Admin" {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}

	employeeID := c.Param("id")
	if err := h.service.DisableEmployee(c.Request.Context(), employeeID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "employee status toggled successfully")
}
