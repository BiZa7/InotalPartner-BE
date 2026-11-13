package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutIntroDocumentHandler struct {
	svc *service.AboutIntroUpsertService
}

func NewAboutIntroDocumentHandler(svc *service.AboutIntroUpsertService) *AboutIntroDocumentHandler {
	return &AboutIntroDocumentHandler{svc: svc}
}

// GetDocumentAdmin — GET /api/about/intro (Admin & Super Admin)
func (h *AboutIntroDocumentHandler) GetDocumentAdmin(c *gin.Context) {
	doc, err := h.svc.GetDocument()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data intro & misi visi",
		"data":    doc,
	})
}

// GetDocumentPublic — GET /api/public/about/intro (Public, No Auth)
func (h *AboutIntroDocumentHandler) GetDocumentPublic(c *gin.Context) {
	doc, err := h.svc.GetDocument()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data intro & misi visi aktif",
		"data":    doc,
	})
}

// SaveDocument — PUT /api/about/intro (Admin & Super Admin)
func (h *AboutIntroDocumentHandler) SaveDocument(c *gin.Context) {
	var req service.AboutIntroDocument
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	if err := h.svc.SaveDocument(req); err != nil {
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
			"message": "Terjadi kesalahan pada server: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Intro & Visi Misi berhasil disimpan",
	})
}

