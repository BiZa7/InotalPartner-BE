package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutVideoHandler struct {
	videoSvc *service.AboutVideoService
}

func NewAboutVideoHandler(videoSvc *service.AboutVideoService) *AboutVideoHandler {
	return &AboutVideoHandler{videoSvc: videoSvc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicVideo — GET /api/public/about/video (Public, No Auth)
// Mengambil daftar video profil aktif untuk halaman publik Tentang Kami.
func (h *AboutVideoHandler) GetPublicVideo(c *gin.Context) {
	videos, err := h.videoSvc.GetPublicVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil video profil publik",
		"data":    videos,
	})
}

// ─── API Pengelolaan Admin ───────────────────────────────────────────────────

// GetVideos — GET /api/about/video (Admin Only)
// Mengambil seluruh video profil untuk dashboard admin (termasuk yang nonaktif).
func (h *AboutVideoHandler) GetVideos(c *gin.Context) {
	videos, err := h.videoSvc.GetAllVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar video profil",
		"data":    videos,
	})
}

// GetVideoByID — GET /api/about/video/:id (Admin Only)
// Mengambil detail video profil berdasarkan ID.
func (h *AboutVideoHandler) GetVideoByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	video, err := h.videoSvc.GetVideoByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "video profil tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail video profil",
		"data":    video,
	})
}

// CreateVideo — POST /api/about/video (Admin Only)
// Membuat record video profil baru.
func (h *AboutVideoHandler) CreateVideo(c *gin.Context) {
	var req service.CreateAboutVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	video, err := h.videoSvc.CreateVideo(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Video profil berhasil ditambahkan",
		"data":    video,
	})
}

// UpdateVideo — PUT /api/about/video/:id (Admin Only)
// Memperbarui video profil berdasarkan ID.
func (h *AboutVideoHandler) UpdateVideo(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateAboutVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	video, err := h.videoSvc.UpdateVideo(id, req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "video profil tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Video profil berhasil diperbarui",
		"data":    video,
	})
}

// DeleteVideo — DELETE /api/about/video/:id (Admin Only)
// Menghapus video profil berdasarkan ID.
func (h *AboutVideoHandler) DeleteVideo(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.videoSvc.DeleteVideo(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "video profil tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Video profil berhasil dihapus",
	})
}
