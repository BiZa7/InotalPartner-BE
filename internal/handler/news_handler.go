package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/model"
	"inotal-be/internal/service"
)

type NewsHandler struct {
	articleSvc *service.ArticleService
}

func NewNewsHandler(articleSvc *service.ArticleService) *NewsHandler {
	return &NewsHandler{articleSvc: articleSvc}
}

// GetNews — GET /api/news?page=1&limit=10
func (h *NewsHandler) GetNews(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi diperlukan",
		})
		return
	}
	userID := userIDVal.(uint)

	userRoleVal, _ := c.Get("user_role")
	userRole, _ := userRoleVal.(string)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.articleSvc.GetArticlesByPostType(userID, userRole, model.PostTypeNews, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar berita",
		"data":    result,
	})
}

// GetNewsByID — GET /api/news/:id
func (h *NewsHandler) GetNewsByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi diperlukan",
		})
		return
	}
	userID := userIDVal.(uint)

	userRoleVal, _ := c.Get("user_role")
	userRole, _ := userRoleVal.(string)

	article, err := h.articleSvc.GetArticleByIDAndPostType(id, userID, userRole, model.PostTypeNews)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "artikel tidak ditemukan":
			status = http.StatusNotFound
		case "akses ditolak: bukan pemilik artikel":
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data berita",
		"data":    article,
	})
}

// CreateNews — POST /api/news
func (h *NewsHandler) CreateNews(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi diperlukan",
		})
		return
	}
	userID := userIDVal.(uint)

	var req service.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	article, err := h.articleSvc.CreateArticleWithPostType(userID, req, model.PostTypeNews)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "kategori tidak ditemukan", "satu atau lebih kategori tidak ditemukan", "satu atau lebih tag tidak ditemukan", "judul artikel tidak boleh kosong", "konten artikel tidak boleh kosong", "kategori wajib dipilih", "judul artikel tidak valid untuk dijadikan slug":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Berita berhasil dibuat",
		"data":    article,
	})
}

// UpdateNews — PUT /api/news/:id
func (h *NewsHandler) UpdateNews(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi diperlukan",
		})
		return
	}
	userID := userIDVal.(uint)

	var req service.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	article, err := h.articleSvc.UpdateArticleWithPostType(id, userID, req, model.PostTypeNews)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "artikel tidak ditemukan":
			status = http.StatusNotFound
		case "akses ditolak: bukan pemilik artikel":
			status = http.StatusForbidden
		case "kategori tidak ditemukan", "satu atau lebih kategori tidak ditemukan", "satu atau lebih tag tidak ditemukan", "judul artikel tidak valid untuk dijadikan slug":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berita berhasil diperbarui",
		"data":    article,
	})
}

// PublishNews — POST /api/news/:id/publish
func (h *NewsHandler) PublishNews(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	article, err := h.articleSvc.PublishArticleWithPostType(id, model.PostTypeNews)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "artikel tidak ditemukan":
			status = http.StatusNotFound
		case "artikel sudah dipublish", "artikel yang sudah diturunkan tidak dapat dipublish kembali":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berita berhasil dipublish",
		"data":    article,
	})
}

// TakedownNews — POST /api/news/:id/takedown
func (h *NewsHandler) TakedownNews(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi diperlukan",
		})
		return
	}
	userID := userIDVal.(uint)

	article, err := h.articleSvc.TakedownArticleWithPostType(id, userID, model.PostTypeNews)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "artikel tidak ditemukan":
			status = http.StatusNotFound
		case "artikel sudah diturunkan", "hanya artikel yang berstatus published yang dapat diturunkan":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berita berhasil diturunkan (takedown)",
		"data":    article,
	})
}
