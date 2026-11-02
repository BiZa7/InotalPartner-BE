package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutVideoHighlightHandler struct {
	highlightSvc *service.AboutVideoHighlightService
}

func NewAboutVideoHighlightHandler(highlightSvc *service.AboutVideoHighlightService) *AboutVideoHighlightHandler {
	return &AboutVideoHighlightHandler{highlightSvc: highlightSvc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicHighlights — GET /api/public/about/video-highlights (Public, No Auth)
// Mengambil daftar highlight aktif saja untuk halaman publik Tentang Kami.
func (h *AboutVideoHighlightHandler) GetPublicHighlights(c *gin.Context) {
	highlights, err := h.highlightSvc.GetPublicHighlights()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar highlight video publik",
		"data":    highlights,
	})
}

// ─── API Pengelolaan Admin ───────────────────────────────────────────────────

// GetHighlights — GET /api/about/video-highlights (Admin & Super Admin)
// Mengambil seluruh highlight video untuk dashboard admin.
func (h *AboutVideoHighlightHandler) GetHighlights(c *gin.Context) {
	highlights, err := h.highlightSvc.GetAllHighlights()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar highlight video",
		"data":    highlights,
	})
}

// GetHighlightByID — GET /api/about/video-highlights/:id (Admin & Super Admin)
// Mengambil detail highlight video berdasarkan ID.
func (h *AboutVideoHighlightHandler) GetHighlightByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	highlight, err := h.highlightSvc.GetHighlightByID(id)
	if err != nil {
		if errors.Is(err, service.ErrAboutVideoHighlightNotFound) {
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
		"message": "Berhasil mengambil detail highlight video",
		"data":    highlight,
	})
}

// CreateHighlight — POST /api/about/video-highlights (Admin & Super Admin)
// Membuat highlight video baru.
func (h *AboutVideoHighlightHandler) CreateHighlight(c *gin.Context) {
	var req service.CreateAboutVideoHighlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	highlight, err := h.highlightSvc.CreateHighlight(req)
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
		"message": "Highlight video berhasil dibuat",
		"data":    highlight,
	})
}

// UpdateHighlight — PUT /api/about/video-highlights/:id (Admin & Super Admin)
// Memperbarui highlight video berdasarkan ID.
func (h *AboutVideoHighlightHandler) UpdateHighlight(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateAboutVideoHighlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	highlight, err := h.highlightSvc.UpdateHighlight(id, req)
	if err != nil {
		if errors.Is(err, service.ErrAboutVideoHighlightNotFound) {
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
		"message": "Highlight video berhasil diperbarui",
		"data":    highlight,
	})
}

// DeleteHighlight — DELETE /api/about/video-highlights/:id (Admin & Super Admin)
// Menghapus highlight video berdasarkan ID.
func (h *AboutVideoHighlightHandler) DeleteHighlight(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.highlightSvc.DeleteHighlight(id); err != nil {
		if errors.Is(err, service.ErrAboutVideoHighlightNotFound) {
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
		"message": "Highlight video berhasil dihapus",
	})
}

