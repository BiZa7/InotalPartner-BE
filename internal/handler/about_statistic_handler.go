package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutStatisticHandler struct {
	statisticSvc *service.AboutStatisticService
}

func NewAboutStatisticHandler(statisticSvc *service.AboutStatisticService) *AboutStatisticHandler {
	return &AboutStatisticHandler{statisticSvc: statisticSvc}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicStatistics — GET /api/public/about/statistics (Public, No Auth)
// Mengambil daftar statistik aktif saja untuk ditampilkan pada halaman publik Tentang Kami.
func (h *AboutStatisticHandler) GetPublicStatistics(c *gin.Context) {
	statistics, err := h.statisticSvc.GetPublicStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar statistik publik",
		"data":    statistics,
	})
}

// ─── API Pengelolaan Admin ───────────────────────────────────────────────────

// GetAllStatistics — GET /api/about/statistics (Admin & Super Admin)
// Mengambil seluruh data statistik untuk dashboard admin (termasuk yang nonaktif).
func (h *AboutStatisticHandler) GetAllStatistics(c *gin.Context) {
	statistics, err := h.statisticSvc.GetAllStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar statistik",
		"data":    statistics,
	})
}

// GetStatisticByID — GET /api/about/statistics/:id (Admin & Super Admin)
// Mengambil detail data statistik berdasarkan ID.
func (h *AboutStatisticHandler) GetStatisticByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	statistic, err := h.statisticSvc.GetStatisticByID(id)
	if err != nil {
		if errors.Is(err, service.ErrAboutStatisticNotFound) {
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
		"message": "Berhasil mengambil detail statistik",
		"data":    statistic,
	})
}

// CreateStatistic — POST /api/about/statistics (Admin & Super Admin)
// Membuat data statistik baru.
func (h *AboutStatisticHandler) CreateStatistic(c *gin.Context) {
	var req service.CreateAboutStatisticRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	statistic, err := h.statisticSvc.CreateStatistic(req)
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
		"message": "Statistik berhasil dibuat",
		"data":    statistic,
	})
}

// UpdateStatistic — PUT /api/about/statistics/:id (Admin & Super Admin)
// Memperbarui data statistik berdasarkan ID.
func (h *AboutStatisticHandler) UpdateStatistic(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateAboutStatisticRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	statistic, err := h.statisticSvc.UpdateStatistic(id, req)
	if err != nil {
		if errors.Is(err, service.ErrAboutStatisticNotFound) {
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
		"message": "Statistik berhasil diperbarui",
		"data":    statistic,
	})
}

// DeleteStatistic — DELETE /api/about/statistics/:id (Admin & Super Admin)
// Menghapus data statistik berdasarkan ID.
func (h *AboutStatisticHandler) DeleteStatistic(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.statisticSvc.DeleteStatistic(id); err != nil {
		if errors.Is(err, service.ErrAboutStatisticNotFound) {
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
		"message": "Statistik berhasil dihapus",
	})
}

