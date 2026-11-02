package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutVideoSettingHandler struct {
	settingSvc *service.AboutVideoSettingService
}

func NewAboutVideoSettingHandler(settingSvc *service.AboutVideoSettingService) *AboutVideoSettingHandler {
	return &AboutVideoSettingHandler{settingSvc: settingSvc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicSetting — GET /api/public/about/video-setting (Public, No Auth)
// Mengambil pengaturan section video profil yang aktif untuk halaman publik.
func (h *AboutVideoSettingHandler) GetPublicSetting(c *gin.Context) {
	setting, err := h.settingSvc.GetPublicSetting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil pengaturan section video profil",
		"data":    setting,
	})
}

// ─── API Pengelolaan Admin ───────────────────────────────────────────────────

// GetSetting — GET /api/about/video-setting (Admin & Super Admin)
// Mengambil pengaturan section video profil untuk dashboard admin.
func (h *AboutVideoSettingHandler) GetSetting(c *gin.Context) {
	setting, err := h.settingSvc.GetSetting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil pengaturan section video profil",
		"data":    setting,
	})
}

// SaveSetting — PUT /api/about/video-setting (Admin & Super Admin)
// Menyimpan atau memperbarui pengaturan section video profil (upsert single-record).
func (h *AboutVideoSettingHandler) SaveSetting(c *gin.Context) {
	var req service.SaveAboutVideoSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	setting, err := h.settingSvc.SaveSetting(req)
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengaturan section video profil berhasil disimpan",
		"data":    setting,
	})
}

