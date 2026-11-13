package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutIntroMetricHandler struct {
	svc *service.AboutIntroMetricService
}

func NewAboutIntroMetricHandler(svc *service.AboutIntroMetricService) *AboutIntroMetricHandler {
	return &AboutIntroMetricHandler{svc: svc}
}

// GetAll — GET /api/about/intro/metrics (Admin & Super Admin)
func (h *AboutIntroMetricHandler) GetAll(c *gin.Context) {
	metrics, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar metric perkenalan tentang kami",
		"data":    metrics,
	})
}

// GetByID — GET /api/about/intro/metrics/:id (Admin & Super Admin)
func (h *AboutIntroMetricHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	metric, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrAboutIntroMetricNotFound) {
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
		"message": "Berhasil mengambil detail metric perkenalan tentang kami",
		"data":    metric,
	})
}

// GetBySettingID — GET /api/about/intro/settings/:setting_id/metrics (Admin & Super Admin)
func (h *AboutIntroMetricHandler) GetBySettingID(c *gin.Context) {
	settingID, err := parseUintParam(c, "setting_id")
	if err != nil {
		return
	}

	metrics, err := h.svc.GetBySettingID(settingID)
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar metric berdasarkan setting ID",
		"data":    metrics,
	})
}

// GetActiveBySettingID — GET /api/about/intro/settings/:setting_id/metrics/active (Admin & Super Admin)
func (h *AboutIntroMetricHandler) GetActiveBySettingID(c *gin.Context) {
	settingID, err := parseUintParam(c, "setting_id")
	if err != nil {
		return
	}

	metrics, err := h.svc.GetActiveBySettingID(settingID)
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar metric aktif berdasarkan setting ID",
		"data":    metrics,
	})
}

// Create — POST /api/about/intro/metrics (Admin & Super Admin)
func (h *AboutIntroMetricHandler) Create(c *gin.Context) {
	var req service.CreateAboutIntroMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	metric, err := h.svc.Create(req)
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
		"message": "Metric perkenalan tentang kami berhasil ditambahkan",
		"data":    metric,
	})
}

// Update — PUT /api/about/intro/metrics/:id (Admin & Super Admin)
func (h *AboutIntroMetricHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateAboutIntroMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	metric, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrAboutIntroMetricNotFound) {
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
		"message": "Metric perkenalan tentang kami berhasil diperbarui",
		"data":    metric,
	})
}

// Delete — DELETE /api/about/intro/metrics/:id (Admin & Super Admin)
func (h *AboutIntroMetricHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrAboutIntroMetricNotFound) {
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
		"message": "Metric perkenalan tentang kami berhasil dihapus",
	})
}

