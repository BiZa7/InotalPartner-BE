package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type HomepageBannerHandler struct {
	bannerSvc *service.HomepageBannerService
}

func NewHomepageBannerHandler(bannerSvc *service.HomepageBannerService) *HomepageBannerHandler {
	return &HomepageBannerHandler{bannerSvc: bannerSvc}
}

// GetBanners — GET /api/homepage/banners (Admin Only)
func (h *HomepageBannerHandler) GetBanners(c *gin.Context) {
	banners, err := h.bannerSvc.GetAllBanners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar banner",
		"data":    banners,
	})
}

// GetBannerByID — GET /api/homepage/banners/:id (Admin Only)
func (h *HomepageBannerHandler) GetBannerByID(c *gin.Context) {
	id, err := parseUintBannerParam(c, "id")
	if err != nil {
		return
	}

	banner, err := h.bannerSvc.GetBannerByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "banner tidak ditemukan" {
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
		"message": "Berhasil mengambil detail banner",
		"data":    banner,
	})
}

// CreateBanner — POST /api/homepage/banners (Admin Only)
func (h *HomepageBannerHandler) CreateBanner(c *gin.Context) {
	var req service.CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data request tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	banner, err := h.bannerSvc.CreateBanner(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Banner berhasil dibuat",
		"data":    banner,
	})
}

// UpdateBanner — PUT /api/homepage/banners/:id (Admin Only)
func (h *HomepageBannerHandler) UpdateBanner(c *gin.Context) {
	id, err := parseUintBannerParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data request tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	banner, err := h.bannerSvc.UpdateBanner(id, req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "banner tidak ditemukan" {
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
		"message": "Banner berhasil diperbarui",
		"data":    banner,
	})
}

// DeleteBanner — DELETE /api/homepage/banners/:id (Admin Only)
func (h *HomepageBannerHandler) DeleteBanner(c *gin.Context) {
	id, err := parseUintBannerParam(c, "id")
	if err != nil {
		return
	}

	if err := h.bannerSvc.DeleteBanner(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "banner tidak ditemukan" {
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
		"message": "Banner berhasil dihapus",
	})
}

// GetPublicBanners — GET /api/public/homepage/banners (Public)
func (h *HomepageBannerHandler) GetPublicBanners(c *gin.Context) {
	banners, err := h.bannerSvc.GetActiveBanners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar banner aktif",
		"data":    banners,
	})
}

// Helper lokal untuk parse uint param URL
func parseUintBannerParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return 0, err
	}
	return uint(val), nil
}
