package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/model"
	"inotal-be/internal/service"
)

type PageBannerHandler struct {
	bannerSvc *service.PageBannerService
}

func NewPageBannerHandler(bannerSvc *service.PageBannerService) *PageBannerHandler {
	return &PageBannerHandler{bannerSvc: bannerSvc}
}

// bannerToResponse transforms model.PageBanner into response format compatible with both BE and FE contracts.
// FE expects: page, label, title, description, image_url.
func bannerToResponse(banner *model.PageBanner, requestedPage string) gin.H {
	if banner == nil {
		return nil
	}

	reverseMap := map[string]string{
		"about":   "tentang",
		"news":    "berita",
		"article": "artikel",
		"contact": "kontak",
		"event":   "event",
	}

	fePageKey := reverseMap[banner.PageType]
	if fePageKey == "" {
		fePageKey = banner.PageType
	}
	// If requested page in FE format, preserve it
	if requestedPage != "" {
		aliasMap := map[string]string{
			"tentang": "about",
			"berita":  "news",
			"artikel": "article",
			"kontak":  "contact",
			"event":   "event",
		}
		if aliasMap[requestedPage] == banner.PageType {
			fePageKey = requestedPage
		}
	}

	label := banner.Subtitle
	subtitle := banner.Subtitle

	return gin.H{
		"id":          banner.ID,
		"page_type":   banner.PageType,
		"page":        fePageKey,
		"image_url":   banner.ImageURL,
		"subtitle":    subtitle,
		"label":       label,
		"title":       banner.Title,
		"description": banner.Description,
		"is_active":   banner.IsActive,
		"created_at":  banner.CreatedAt,
		"updated_at":  banner.UpdatedAt,
	}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicBanner — GET /api/public/page-banners/:page_type (Public, No Auth)
// Mengambil banner aktif saja untuk pengunjung halaman publik.
func (h *PageBannerHandler) GetPublicBanner(c *gin.Context) {
	pageType := c.Param("page_type")

	banner, err := h.bannerSvc.GetPublicBannerByPageType(pageType)
	if err != nil {
		if errors.Is(err, service.ErrPageBannerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Banner untuk halaman ini belum tersedia atau sedang nonaktif",
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
		"message": "Berhasil mengambil data banner publik",
		"data":    bannerToResponse(banner, pageType),
	})
}

// ─── API Pengelolaan Admin ───────────────────────────────────────────────────

// GetBanner — GET /api/page-banners/:page_type (Admin & Super Admin)
// Mengambil data banner berdasarkan page_type untuk dashboard admin.
// Bila belum ada, kembalikan HTTP 200 {"success": true, "data": null} sesuai kontrak frontend admin.
func (h *PageBannerHandler) GetBanner(c *gin.Context) {
	pageType := c.Param("page_type")

	banner, err := h.bannerSvc.GetBannerByPageType(pageType)
	if err != nil {
		if errors.Is(err, service.ErrPageBannerNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Banner belum tersedia",
				"data":    nil,
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
		"message": "Berhasil mengambil data banner",
		"data":    bannerToResponse(banner, pageType),
	})
}

// CreateBanner — POST /api/page-banners/:page_type (Admin & Super Admin)
// Membuat banner baru untuk page_type tertentu.
func (h *PageBannerHandler) CreateBanner(c *gin.Context) {
	pageType := c.Param("page_type")

	var req service.CreatePageBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	banner, err := h.bannerSvc.CreateBanner(pageType, req)
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
		"message": "Banner berhasil dibuat",
		"data":    bannerToResponse(banner, pageType),
	})
}

// UpdateBanner — PUT /api/page-banners/:page_type (Admin & Super Admin)
// Memperbarui data banner untuk page_type tertentu.
func (h *PageBannerHandler) UpdateBanner(c *gin.Context) {
	pageType := c.Param("page_type")

	var req service.UpdatePageBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
		})
		return
	}

	banner, err := h.bannerSvc.UpdateBanner(pageType, req)
	if err != nil {
		if errors.Is(err, service.ErrPageBannerNotFound) {
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
		"message": "Banner berhasil diperbarui",
		"data":    bannerToResponse(banner, pageType),
	})
}

// DeleteBanner — DELETE /api/page-banners/:page_type (Admin & Super Admin)
// Menghapus banner untuk page_type tertentu.
func (h *PageBannerHandler) DeleteBanner(c *gin.Context) {
	pageType := c.Param("page_type")

	if err := h.bannerSvc.DeleteBanner(pageType); err != nil {
		if errors.Is(err, service.ErrPageBannerNotFound) {
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
		"message": "Banner berhasil dihapus",
	})
}

