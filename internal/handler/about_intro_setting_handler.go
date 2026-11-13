package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutIntroSettingHandler struct {
	svc *service.AboutIntroSettingService
}

func NewAboutIntroSettingHandler(svc *service.AboutIntroSettingService) *AboutIntroSettingHandler {
	return &AboutIntroSettingHandler{svc: svc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetActivePublic — GET /api/public/about/intro (Public, No Auth)
func (h *AboutIntroSettingHandler) GetActivePublic(c *gin.Context) {
	resp, err := h.svc.GetActivePublicResponse()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	if resp == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Belum ada konfigurasi perkenalan tentang kami yang aktif",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil konfigurasi perkenalan tentang kami aktif",
		"data":    resp,
	})
}

// ─── API Admin ────────────────────────────────────────────────────────────────

// GetAll — GET /api/about/intro/settings (Admin & Super Admin)
func (h *AboutIntroSettingHandler) GetAll(c *gin.Context) {
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
		"message": "Berhasil mengambil seluruh konfigurasi perkenalan tentang kami",
		"data":    settings,
	})
}

// GetByID — GET /api/about/intro/settings/:id (Admin & Super Admin)
func (h *AboutIntroSettingHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	setting, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrAboutIntroSettingNotFound) {
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
		"message": "Berhasil mengambil detail konfigurasi perkenalan tentang kami",
		"data":    setting,
	})
}

// Create — POST /api/about/intro/settings (Admin & Super Admin)
func (h *AboutIntroSettingHandler) Create(c *gin.Context) {
	var req service.CreateAboutIntroSettingRequest
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
		"message": "Konfigurasi perkenalan tentang kami berhasil dibuat",
		"data":    setting,
	})
}

// Update — PUT /api/about/intro/settings/:id (Admin & Super Admin)
func (h *AboutIntroSettingHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateAboutIntroSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	setting, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrAboutIntroSettingNotFound) {
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
		"message": "Konfigurasi perkenalan tentang kami berhasil diperbarui",
		"data":    setting,
	})
}

// Delete — DELETE /api/about/intro/settings/:id (Admin & Super Admin)
func (h *AboutIntroSettingHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrAboutIntroSettingNotFound) {
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
		"message": "Konfigurasi perkenalan tentang kami berhasil dihapus",
	})
}

