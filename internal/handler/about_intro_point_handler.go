package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutIntroPointHandler struct {
	svc *service.AboutIntroPointService
}

func NewAboutIntroPointHandler(svc *service.AboutIntroPointService) *AboutIntroPointHandler {
	return &AboutIntroPointHandler{svc: svc}
}

// GetAll — GET /api/about/intro/points (Admin & Super Admin)
func (h *AboutIntroPointHandler) GetAll(c *gin.Context) {
	points, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar point perkenalan tentang kami",
		"data":    points,
	})
}

// GetByID — GET /api/about/intro/points/:id (Admin & Super Admin)
func (h *AboutIntroPointHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	point, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrAboutIntroPointNotFound) {
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
		"message": "Berhasil mengambil detail point perkenalan tentang kami",
		"data":    point,
	})
}

// GetBySettingID — GET /api/about/intro/settings/:setting_id/points (Admin & Super Admin)
func (h *AboutIntroPointHandler) GetBySettingID(c *gin.Context) {
	settingID, err := parseUintParam(c, "setting_id")
	if err != nil {
		return
	}

	points, err := h.svc.GetBySettingID(settingID)
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
		"message": "Berhasil mengambil daftar point berdasarkan setting ID",
		"data":    points,
	})
}

// GetBySettingIDAndType — GET /api/about/intro/settings/:setting_id/points/:section_type (Admin & Super Admin)
func (h *AboutIntroPointHandler) GetBySettingIDAndType(c *gin.Context) {
	settingID, err := parseUintParam(c, "setting_id")
	if err != nil {
		return
	}

	sectionType := c.Param("section_type")

	points, err := h.svc.GetBySettingIDAndType(settingID, sectionType)
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
		"message": "Berhasil mengambil daftar point berdasarkan setting ID dan tipe seksi",
		"data":    points,
	})
}

// GetActiveBySettingID — GET /api/about/intro/settings/:setting_id/points/active (Admin & Super Admin)
func (h *AboutIntroPointHandler) GetActiveBySettingID(c *gin.Context) {
	settingID, err := parseUintParam(c, "setting_id")
	if err != nil {
		return
	}

	points, err := h.svc.GetActiveBySettingID(settingID)
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
		"message": "Berhasil mengambil daftar point aktif berdasarkan setting ID",
		"data":    points,
	})
}

// Create — POST /api/about/intro/points (Admin & Super Admin)
func (h *AboutIntroPointHandler) Create(c *gin.Context) {
	var req service.CreateAboutIntroPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	point, err := h.svc.Create(req)
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
		"message": "Point perkenalan tentang kami berhasil ditambahkan",
		"data":    point,
	})
}

// Update — PUT /api/about/intro/points/:id (Admin & Super Admin)
func (h *AboutIntroPointHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateAboutIntroPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	point, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrAboutIntroPointNotFound) {
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
		"message": "Point perkenalan tentang kami berhasil diperbarui",
		"data":    point,
	})
}

// Delete — DELETE /api/about/intro/points/:id (Admin & Super Admin)
func (h *AboutIntroPointHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrAboutIntroPointNotFound) {
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
		"message": "Point perkenalan tentang kami berhasil dihapus",
	})
}

