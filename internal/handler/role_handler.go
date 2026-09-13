package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type RoleHandler struct {
	roleSvc *service.RoleService
}

func NewRoleHandler(roleSvc *service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

// GetRoles — GET /api/roles
func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleSvc.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar role",
		"data":    roles,
	})
}

// GetRole — GET /api/roles/:id
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	role, err := h.roleSvc.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data role",
		"data":    role,
	})
}

// CreateRole — POST /api/roles
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	role, err := h.roleSvc.CreateRole(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "nama role sudah digunakan" {
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
		"message": "Role berhasil dibuat",
		"data":    role,
	})
}

// UpdateRole — PUT /api/roles/:id
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	role, err := h.roleSvc.UpdateRole(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "role tidak ditemukan":
			status = http.StatusNotFound
		case "nama role sudah digunakan":
			status = http.StatusConflict
		case "nama role system tidak dapat diubah":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role berhasil diperbarui",
		"data":    role,
	})
}

// DeleteRole — DELETE /api/roles/:id
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.roleSvc.DeleteRole(id); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "role tidak ditemukan":
			status = http.StatusNotFound
		case "role system tidak dapat dihapus", "role tidak dapat dihapus karena masih digunakan oleh user":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role berhasil dihapus",
	})
}
