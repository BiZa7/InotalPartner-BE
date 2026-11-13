package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type HomeCTAHandler struct {
	ctaSvc *service.HomeCTAService
}

func NewHomeCTAHandler(ctaSvc *service.HomeCTAService) *HomeCTAHandler {
	return &HomeCTAHandler{ctaSvc: ctaSvc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicCTA — GET /api/public/home/cta (Public, No Auth)
func (h *HomeCTAHandler) GetPublicCTA(c *gin.Context) {
	cta, err := h.ctaSvc.GetActivePublic()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data CTA beranda aktif",
		"data":    cta,
	})
}

// ─── API Admin ────────────────────────────────────────────────────────────────

// GetAllCTA — GET /api/home/cta (Admin & Super Admin)
func (h *HomeCTAHandler) GetAllCTA(c *gin.Context) {
	ctas, err := h.ctaSvc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar CTA beranda",
		"data":    ctas,
	})
}

// GetCTAByID — GET /api/home/cta/:id (Admin & Super Admin)
func (h *HomeCTAHandler) GetCTAByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	cta, err := h.ctaSvc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrHomeCTANotFound) {
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
		"message": "Berhasil mengambil detail CTA beranda",
		"data":    cta,
	})
}

// CreateCTA — POST /api/home/cta (Admin & Super Admin)
func (h *HomeCTAHandler) CreateCTA(c *gin.Context) {
	var req service.CreateHomeCTARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	cta, err := h.ctaSvc.Create(req)
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
			"message": "Gagal membuat CTA beranda: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "CTA beranda berhasil dibuat",
		"data":    cta,
	})
}

// UpdateCTA — PUT /api/home/cta/:id (Admin & Super Admin)
func (h *HomeCTAHandler) UpdateCTA(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateHomeCTARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	cta, err := h.ctaSvc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrHomeCTANotFound) {
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
			"message": "Gagal memperbarui CTA beranda: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "CTA beranda berhasil diperbarui",
		"data":    cta,
	})
}

// DeleteCTA — DELETE /api/home/cta/:id (Admin & Super Admin)
func (h *HomeCTAHandler) DeleteCTA(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.ctaSvc.Delete(id); err != nil {
		if errors.Is(err, service.ErrHomeCTANotFound) {
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
			"message": "Gagal menghapus CTA beranda: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "CTA beranda berhasil dihapus",
	})
}

