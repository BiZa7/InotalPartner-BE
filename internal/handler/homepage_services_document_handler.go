package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type HomepageServicesDocumentHandler struct {
	svc *service.HomepageServicesUpsertService
}

func NewHomepageServicesDocumentHandler(svc *service.HomepageServicesUpsertService) *HomepageServicesDocumentHandler {
	return &HomepageServicesDocumentHandler{svc: svc}
}

// GetDocumentAdmin — GET /api/homepage/services (Admin & Super Admin)
func (h *HomepageServicesDocumentHandler) GetDocumentAdmin(c *gin.Context) {
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
		"message": "Berhasil mengambil data layanan beranda",
		"data":    doc,
	})
}

// GetDocumentPublic — GET /api/public/homepage/services (Public, No Auth)
func (h *HomepageServicesDocumentHandler) GetDocumentPublic(c *gin.Context) {
	doc, err := h.svc.GetDocument()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data layanan beranda",
		"data":    doc,
	})
}

// SaveDocument — PUT /api/homepage/services (Admin & Super Admin)
func (h *HomepageServicesDocumentHandler) SaveDocument(c *gin.Context) {
	var req service.HomepageServicesDocument
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
		"message": "Layanan Beranda berhasil disimpan",
	})
}

