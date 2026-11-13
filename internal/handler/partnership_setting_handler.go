package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type PartnershipSettingHandler struct {
	svc *service.PartnershipSettingService
}

func NewPartnershipSettingHandler(svc *service.PartnershipSettingService) *PartnershipSettingHandler {
	return &PartnershipSettingHandler{svc: svc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetActivePublic — GET /api/public/home/partnership/settings (Public, No Auth)
func (h *PartnershipSettingHandler) GetActivePublic(c *gin.Context) {
	setting, err := h.svc.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil konfigurasi partnership beranda aktif",
		"data":    setting,
	})
}

// ─── API Admin ────────────────────────────────────────────────────────────────

// GetAll — GET /api/home/partnership/settings (Admin & Super Admin)
func (h *PartnershipSettingHandler) GetAll(c *gin.Context) {
	settings, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh konfigurasi partnership beranda",
		"data":    settings,
	})
}

// GetByID — GET /api/home/partnership/settings/:id (Admin & Super Admin)
func (h *PartnershipSettingHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	setting, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrPartnershipSettingNotFound) {
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
		"message": "Berhasil mengambil detail konfigurasi partnership beranda",
		"data":    setting,
	})
}

// Create — POST /api/home/partnership/settings (Admin & Super Admin)
func (h *PartnershipSettingHandler) Create(c *gin.Context) {
	var req service.CreatePartnershipSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	setting, err := h.svc.Create(req)
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
		"message": "Konfigurasi partnership beranda berhasil dibuat",
		"data":    setting,
	})
}

// Update — PUT /api/home/partnership/settings/:id (Admin & Super Admin)
func (h *PartnershipSettingHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdatePartnershipSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	setting, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrPartnershipSettingNotFound) {
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
		"message": "Konfigurasi partnership beranda berhasil diperbarui",
		"data":    setting,
	})
}

// Delete — DELETE /api/home/partnership/settings/:id (Admin & Super Admin)
func (h *PartnershipSettingHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrPartnershipSettingNotFound) {
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
		"message": "Konfigurasi partnership beranda berhasil dihapus",
	})
}

