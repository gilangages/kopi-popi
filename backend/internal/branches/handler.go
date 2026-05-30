package branches

import (
	"strconv"

	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// GetAllBranches (Public / Admin)
func (h *Handler) GetAllBranches(c *gin.Context) {
	// Ambil role dari JWT jika ada (ini untuk rute public yang mungkin mengirimkan token)
	// Jika rute ini public murni tanpa JWT middleware, kita perlu handle gracefully.
	// Di main.go, lebih baik route /branches tidak ditaruh di dalam RequireAuth,
	// namun kita bisa mengekstrak token manual jika mau, atau pisahkan.
	// Kita asumsikan token dimasukkan secara opsional (lihat main.go nanti).
	role := c.GetString("role") 
	
	// Cek query parameter `include_inactive`
	includeInactiveStr := c.Query("include_inactive")
	includeInactive := includeInactiveStr == "true"
	
	branches, err := h.service.GetAllBranches(c.Request.Context(), role, includeInactive)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	
	response.Success(c, 200, branches)
}

// CreateBranch (Admin Only)
func (h *Handler) CreateBranch(c *gin.Context) {
	role := c.GetString("role")
	if role != "Admin" {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	
	var req CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	
	if err := h.service.CreateBranch(c.Request.Context(), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	
	response.Success(c, 200, "branch created successfully")
}

// UpdateBranch (Admin Only)
func (h *Handler) UpdateBranch(c *gin.Context) {
	role := c.GetString("role")
	if role != "Admin" {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, 400, "invalid branch id")
		return
	}
	
	var req UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	
	if err := h.service.UpdateBranch(c.Request.Context(), id, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	
	response.Success(c, 200, "branch updated successfully")
}

// DeleteBranch (Admin Only)
func (h *Handler) DeleteBranch(c *gin.Context) {
	role := c.GetString("role")
	if role != "Admin" {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, 400, "invalid branch id")
		return
	}
	
	if err := h.service.DeleteBranch(c.Request.Context(), id); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	
	response.Success(c, 200, "branch soft-deleted successfully")
}
