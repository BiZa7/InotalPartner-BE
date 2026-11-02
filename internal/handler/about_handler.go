package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutHandler struct {
	settingSvc *service.AboutPageSettingService
	memberSvc  *service.TeamMemberService
}

func NewAboutHandler(
	settingSvc *service.AboutPageSettingService,
	memberSvc *service.TeamMemberService,
) *AboutHandler {
	return &AboutHandler{
		settingSvc: settingSvc,
		memberSvc:  memberSvc,
	}
}

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type SaveAboutSettingRequest struct {
	Title     string `json:"title"`
	TeamTitle string `json:"team_title"`
	IsActive  *bool  `json:"is_active"`
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicAboutSetting — GET /api/public/about (Public, No Auth)
// Mengambil judul halaman Tentang yang sedang aktif.
// Jika belum ada konfigurasi aktif, mengembalikan HTTP 200 dengan data null tanpa membuat data default.
func (h *AboutHandler) GetPublicAboutSetting(c *gin.Context) {
	setting, err := h.settingSvc.GetActiveSetting()
	if err != nil {
		if err.Error() == "pengaturan judul aktif tidak ditemukan" {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Belum ada konfigurasi halaman Tentang yang aktif",
				"data":    nil,
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
		"message": "Berhasil mengambil konfigurasi halaman Tentang",
		"data":    setting,
	})
}

// GetPublicTeamMembers — GET /api/public/about/team (Public, No Auth)
// Mengambil daftar anggota tim aktif yang memiliki nama dan foto valid,
// diurutkan berdasarkan display_order ASC, id ASC.
func (h *AboutHandler) GetPublicTeamMembers(c *gin.Context) {
	members, err := h.memberSvc.GetPublicTeamMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar anggota tim publik",
		"data":    members,
	})
}

// ─── API Pengelolaan Admin (Pengaturan Judul) ─────────────────────────────────

// GetActiveSetting — GET /api/about (Admin Only)
// Mengambil konfigurasi judul aktif.
func (h *AboutHandler) GetActiveSetting(c *gin.Context) {
	setting, err := h.settingSvc.GetActiveSetting()
	if err != nil {
		if err.Error() == "pengaturan judul aktif tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Pengaturan judul aktif tidak ditemukan",
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
		"message": "Berhasil mengambil konfigurasi judul",
		"data":    setting,
	})
}

// SaveActiveSetting — PUT /api/about (Admin Only)
// Membuat atau memperbarui judul aktif (upsert) tanpa field section_label.
func (h *AboutHandler) SaveActiveSetting(c *gin.Context) {
	var req SaveAboutSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	setting, err := h.settingSvc.SaveActiveTitle(req.Title, req.TeamTitle, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Konfigurasi judul berhasil disimpan",
		"data":    setting,
	})
}

// ─── API Pengelolaan Admin (Anggota Tim) ──────────────────────────────────────

// GetTeamMembers — GET /api/about/team (Admin Only)
// Mengambil semua anggota tim, termasuk yang nonaktif dan belum lengkap.
func (h *AboutHandler) GetTeamMembers(c *gin.Context) {
	members, err := h.memberSvc.GetAllTeamMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar anggota tim",
		"data":    members,
	})
}

// CreateTeamMember — POST /api/about/team (Admin Only)
// Menambahkan anggota tim baru.
func (h *AboutHandler) CreateTeamMember(c *gin.Context) {
	var req service.CreateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	member, err := h.memberSvc.CreateTeamMember(req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak boleh bernilai negatif") {
			status = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "display_order sudah digunakan") {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Anggota tim berhasil ditambahkan",
		"data":    member,
	})
}

// GetTeamMemberByID — GET /api/about/team/:id (Admin Only)
// Mengambil detail satu anggota tim berdasarkan ID.
func (h *AboutHandler) GetTeamMemberByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	member, err := h.memberSvc.GetTeamMemberByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "anggota tim tidak ditemukan" {
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
		"message": "Berhasil mengambil detail anggota tim",
		"data":    member,
	})
}

// UpdateTeamMember — PUT /api/about/team/:id (Admin Only)
// Memperbarui data anggota tim berdasarkan ID.
func (h *AboutHandler) UpdateTeamMember(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	member, err := h.memberSvc.UpdateTeamMember(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "anggota tim tidak ditemukan" {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "tidak boleh bernilai negatif") {
			status = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "display_order sudah digunakan") {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Anggota tim berhasil diperbarui",
		"data":    member,
	})
}

// DeleteTeamMember — DELETE /api/about/team/:id (Admin Only)
// Menghapus data anggota tim berdasarkan ID.
func (h *AboutHandler) DeleteTeamMember(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.memberSvc.DeleteTeamMember(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "anggota tim tidak ditemukan" {
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
		"message": "Anggota tim berhasil dihapus",
	})
}
