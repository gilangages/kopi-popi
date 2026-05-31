package catalog

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

func parseID(c *gin.Context) (int, error) {
	idParam := c.Param("id")
	return strconv.Atoi(idParam)
}

func checkAdmin(c *gin.Context) bool {
	return c.GetString("role") == "Admin"
}

func checkAdminOrManager(c *gin.Context) bool {
	role := c.GetString("role")
	return role == "Admin" || role == "Manager"
}

// -- Category Handlers --

func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories(c.Request.Context())
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, categories)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	if err := h.service.CreateCategory(c.Request.Context(), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "category created successfully")
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "invalid category id")
		return
	}
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	if err := h.service.UpdateCategory(c.Request.Context(), id, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "category updated successfully")
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "invalid category id")
		return
	}
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "category deleted successfully")
}

// -- Material Handlers --

func (h *Handler) GetAllMaterials(c *gin.Context) {
	if !checkAdminOrManager(c) {
		response.Error(c, 403, "forbidden: admin or manager role required")
		return
	}
	materials, err := h.service.GetAllMaterials(c.Request.Context())
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, materials)
}

func (h *Handler) CreateMaterial(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	var req MaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	if err := h.service.CreateMaterial(c.Request.Context(), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "material created successfully")
}

func (h *Handler) UpdateMaterial(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "invalid material id")
		return
	}
	var req MaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	if err := h.service.UpdateMaterial(c.Request.Context(), id, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "material updated successfully")
}

func (h *Handler) DeleteMaterial(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "invalid material id")
		return
	}
	if err := h.service.DeleteMaterial(c.Request.Context(), id); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "material deleted successfully")
}

// -- Product Handlers --

func (h *Handler) GetAllProducts(c *gin.Context) {
	var categoryID *int
	if catStr := c.Query("category_id"); catStr != "" {
		cat, err := strconv.Atoi(catStr)
		if err == nil {
			categoryID = &cat
		}
	}
	search := c.Query("search")

	products, err := h.service.GetAllProducts(c.Request.Context(), categoryID, search)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, products)
}

func (h *Handler) GetProductDetail(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "invalid product id")
		return
	}

	role := c.GetString("role")
	includeRecipe := c.Query("include_recipe") == "true"

	product, err := h.service.GetProductDetail(c.Request.Context(), id, role, includeRecipe)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, product)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	if err := h.service.CreateProduct(c.Request.Context(), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "product created successfully")
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "invalid product id")
		return
	}
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	if err := h.service.UpdateProduct(c.Request.Context(), id, req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "product updated successfully")
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	if !checkAdmin(c) {
		response.Error(c, 403, "forbidden: admin role required")
		return
	}
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "invalid product id")
		return
	}
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "product deleted successfully")
}
