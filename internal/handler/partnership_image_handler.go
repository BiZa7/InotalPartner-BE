package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type PartnershipImageHandler struct {
	svc *service.PartnershipImageService
}

func NewPartnershipImageHandler(svc *service.PartnershipImageService) *PartnershipImageHandler {
	return &PartnershipImageHandler{svc: svc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetActivePublic — GET /api/public/home/partnership/images (Public, No Auth)
func (h *PartnershipImageHandler) GetActivePublic(c *gin.Context) {
	images, err := h.svc.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil gambar partnership beranda aktif",
		"data":    images,
	})
}

// ─── API Admin ────────────────────────────────────────────────────────────────

// GetAll — GET /api/home/partnership/images (Admin & Super Admin)
func (h *PartnershipImageHandler) GetAll(c *gin.Context) {
	images, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar gambar partnership beranda",
		"data":    images,
	})
}

// GetByID — GET /api/home/partnership/images/:id (Admin & Super Admin)
func (h *PartnershipImageHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	image, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrPartnershipImageNotFound) {
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
		"message": "Berhasil mengambil detail gambar partnership beranda",
		"data":    image,
	})
}

// Create — POST /api/home/partnership/images (Admin & Super Admin)
func (h *PartnershipImageHandler) Create(c *gin.Context) {
	var req service.CreatePartnershipImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	image, err := h.svc.Create(req)
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
		"message": "Gambar partnership beranda berhasil ditambahkan",
		"data":    image,
	})
}

// Update — PUT /api/home/partnership/images/:id (Admin & Super Admin)
func (h *PartnershipImageHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdatePartnershipImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	image, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrPartnershipImageNotFound) {
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
		"message": "Gambar partnership beranda berhasil diperbarui",
		"data":    image,
	})
}

// Delete — DELETE /api/home/partnership/images/:id (Admin & Super Admin)
func (h *PartnershipImageHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrPartnershipImageNotFound) {
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
		"message": "Gambar partnership beranda berhasil dihapus",
	})
}

