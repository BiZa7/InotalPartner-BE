package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type PartnershipServiceHandler struct {
	svc *service.PartnershipServiceService
}

func NewPartnershipServiceHandler(svc *service.PartnershipServiceService) *PartnershipServiceHandler {
	return &PartnershipServiceHandler{svc: svc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetActivePublic — GET /api/public/home/partnership/services (Public, No Auth)
func (h *PartnershipServiceHandler) GetActivePublic(c *gin.Context) {
	services, err := h.svc.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil layanan partnership beranda aktif",
		"data":    services,
	})
}

// ─── API Admin ────────────────────────────────────────────────────────────────

// GetAll — GET /api/home/partnership/services (Admin & Super Admin)
func (h *PartnershipServiceHandler) GetAll(c *gin.Context) {
	services, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar layanan partnership beranda",
		"data":    services,
	})
}

// GetByID — GET /api/home/partnership/services/:id (Admin & Super Admin)
func (h *PartnershipServiceHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	serviceItem, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrPartnershipServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		var valErr *service.ValidationError
		if errors.As(err, &valErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": valErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail layanan partnership beranda",
		"data":    serviceItem,
	})
}

// Create — POST /api/home/partnership/services (Admin & Super Admin)
func (h *PartnershipServiceHandler) Create(c *gin.Context) {
	var req service.CreatePartnershipServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	serviceItem, err := h.svc.Create(req)
	if err != nil {
		var valErr *service.ValidationError
		if errors.As(err, &valErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": valErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Layanan partnership beranda berhasil ditambahkan",
		"data":    serviceItem,
	})
}

// Update — PUT /api/home/partnership/services/:id (Admin & Super Admin)
func (h *PartnershipServiceHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdatePartnershipServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	serviceItem, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrPartnershipServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		var valErr *service.ValidationError
		if errors.As(err, &valErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": valErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Layanan partnership beranda berhasil diperbarui",
		"data":    serviceItem,
	})
}

// Delete — DELETE /api/home/partnership/services/:id (Admin & Super Admin)
func (h *PartnershipServiceHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrPartnershipServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		var valErr *service.ValidationError
		if errors.As(err, &valErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": valErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Layanan partnership beranda berhasil dihapus",
	})
}

