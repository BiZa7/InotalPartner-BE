package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type CategoryHandler struct {
	categorySvc *service.CategoryService
}

func NewCategoryHandler(categorySvc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc: categorySvc}
}

// GetCategories — GET /api/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categorySvc.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar kategori",
		"data":    categories,
	})
}

// GetCategory — GET /api/categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	category, err := h.categorySvc.GetCategoryByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "kategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data kategori",
		"data":    category,
	})
}

// CreateCategory — POST /api/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	category, err := h.categorySvc.CreateCategory(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "kategori sudah ada" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kategori berhasil dibuat",
		"data":    category,
	})
}

// UpdateCategory — PUT /api/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	category, err := h.categorySvc.UpdateCategory(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "kategori tidak ditemukan":
			status = http.StatusNotFound
		case "kategori sudah ada":
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kategori berhasil diperbarui",
		"data":    category,
	})
}

// DeleteCategory — DELETE /api/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.categorySvc.DeleteCategory(id); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "kategori tidak ditemukan":
			status = http.StatusNotFound
		case "kategori masih digunakan oleh artikel":
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kategori berhasil dihapus",
	})
}
