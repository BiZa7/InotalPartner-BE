package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutFAQHandler struct {
	faqSvc *service.AboutFAQService
}

func NewAboutFAQHandler(faqSvc *service.AboutFAQService) *AboutFAQHandler {
	return &AboutFAQHandler{faqSvc: faqSvc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicFAQs — GET /api/public/about/faqs (Public, No Auth)
// Mengambil daftar FAQ aktif saja untuk halaman publik Tentang Kami.
func (h *AboutFAQHandler) GetPublicFAQs(c *gin.Context) {
	faqs, err := h.faqSvc.GetPublicFAQs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar FAQ publik",
		"data":    faqs,
	})
}

// ─── API Pengelolaan Admin ───────────────────────────────────────────────────

// GetFAQs — GET /api/about/faqs (Admin & Super Admin)
// Mengambil seluruh FAQ untuk dashboard admin.
func (h *AboutFAQHandler) GetFAQs(c *gin.Context) {
	faqs, err := h.faqSvc.GetAllFAQs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar FAQ",
		"data":    faqs,
	})
}

// GetFAQByID — GET /api/about/faqs/:id (Admin & Super Admin)
// Mengambil detail FAQ berdasarkan ID.
func (h *AboutFAQHandler) GetFAQByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	faq, err := h.faqSvc.GetFAQByID(id)
	if err != nil {
		if errors.Is(err, service.ErrAboutFAQNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
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
		"message": "Berhasil mengambil detail FAQ",
		"data":    faq,
	})
}

// CreateFAQ — POST /api/about/faqs (Admin & Super Admin)
// Membuat FAQ baru.
func (h *AboutFAQHandler) CreateFAQ(c *gin.Context) {
	var req service.CreateAboutFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	faq, err := h.faqSvc.CreateFAQ(req)
	if err != nil {
		var valErr *service.ValidationError
		if errors.As(err, &valErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": valErr.Error(),
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
		"message": "FAQ berhasil ditambahkan",
		"data":    faq,
	})
}

// UpdateFAQ — PUT /api/about/faqs/:id (Admin & Super Admin)
// Memperbarui data FAQ berdasarkan ID.
func (h *AboutFAQHandler) UpdateFAQ(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateAboutFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	faq, err := h.faqSvc.UpdateFAQ(id, req)
	if err != nil {
		if errors.Is(err, service.ErrAboutFAQNotFound) {
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
				"message": valErr.Error(),
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
		"message": "FAQ berhasil diperbarui",
		"data":    faq,
	})
}

// DeleteFAQ — DELETE /api/about/faqs/:id (Admin & Super Admin)
// Menghapus FAQ berdasarkan ID.
func (h *AboutFAQHandler) DeleteFAQ(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.faqSvc.DeleteFAQ(id); err != nil {
		if errors.Is(err, service.ErrAboutFAQNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
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
		"message": "FAQ berhasil dihapus",
	})
}

