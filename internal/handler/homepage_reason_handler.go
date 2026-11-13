package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type HomepageReasonHandler struct {
	reasonSvc *service.HomepageReasonService
}

func NewHomepageReasonHandler(reasonSvc *service.HomepageReasonService) *HomepageReasonHandler {
	return &HomepageReasonHandler{reasonSvc: reasonSvc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicReasons — GET /api/public/homepage/reasons (Public, No Auth)
func (h *HomepageReasonHandler) GetPublicReasons(c *gin.Context) {
	data, err := h.reasonSvc.GetPublicReasonsResponse()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data alasan mengapa bergabung",
		"data":    data,
	})
}

// ─── Pengaturan Section (Admin & Super Admin) ─────────────────────────────────

// GetSetting — GET /api/homepage/reason-settings
func (h *HomepageReasonHandler) GetSetting(c *gin.Context) {
	setting, err := h.reasonSvc.GetSetting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil pengaturan section mengapa bergabung",
		"data":    setting,
	})
}

// SaveSetting — PUT /api/homepage/reason-settings
func (h *HomepageReasonHandler) SaveSetting(c *gin.Context) {
	var req service.SaveHomepageReasonSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	setting, err := h.reasonSvc.SaveSetting(req)
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
			"message": "Gagal menyimpan pengaturan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengaturan section mengapa bergabung berhasil disimpan",
		"data":    setting,
	})
}

// ─── CRUD Kartu Alasan (Admin & Super Admin) ──────────────────────────────────

// GetAllReasons — GET /api/homepage/reasons
func (h *HomepageReasonHandler) GetAllReasons(c *gin.Context) {
	reasons, err := h.reasonSvc.GetAllReasons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar kartu alasan",
		"data":    reasons,
	})
}

// GetReasonByID — GET /api/homepage/reasons/:id
func (h *HomepageReasonHandler) GetReasonByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	reason, err := h.reasonSvc.GetReasonByID(id)
	if err != nil {
		if errors.Is(err, service.ErrHomepageReasonNotFound) {
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
		"message": "Berhasil mengambil detail kartu alasan",
		"data":    reason,
	})
}

// CreateReason — POST /api/homepage/reasons
func (h *HomepageReasonHandler) CreateReason(c *gin.Context) {
	var req service.CreateHomepageReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	reason, err := h.reasonSvc.CreateReason(req)
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
			"message": "Gagal menambahkan kartu alasan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kartu alasan berhasil ditambahkan",
		"data":    reason,
	})
}

// UpdateReason — PUT /api/homepage/reasons/:id
func (h *HomepageReasonHandler) UpdateReason(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateHomepageReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	reason, err := h.reasonSvc.UpdateReason(id, req)
	if err != nil {
		if errors.Is(err, service.ErrHomepageReasonNotFound) {
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
			"message": "Gagal memperbarui kartu alasan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kartu alasan berhasil diperbarui",
		"data":    reason,
	})
}

// DeleteReason — DELETE /api/homepage/reasons/:id
func (h *HomepageReasonHandler) DeleteReason(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.reasonSvc.DeleteReason(id); err != nil {
		if errors.Is(err, service.ErrHomepageReasonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus kartu alasan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kartu alasan berhasil dihapus",
	})
}

