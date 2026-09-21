package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type ArticleHandler struct {
	articleSvc *service.ArticleService
}

func NewArticleHandler(articleSvc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleSvc: articleSvc}
}

// GetArticles — GET /api/articles?page=1&limit=10
func (h *ArticleHandler) GetArticles(c *gin.Context) {
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

	result, err := h.articleSvc.GetArticles(userID, userRole, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar artikel",
		"data":    result,
	})
}

// GetArticle — GET /api/articles/:id
func (h *ArticleHandler) GetArticle(c *gin.Context) {
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

	article, err := h.articleSvc.GetArticleByID(id, userID, userRole)
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
		"message": "Berhasil mengambil data artikel",
		"data":    article,
	})
}

// CreateArticle — POST /api/articles
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
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

	article, err := h.articleSvc.CreateArticle(userID, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "kategori tidak ditemukan", "satu atau lebih tag tidak ditemukan", "judul artikel tidak boleh kosong", "konten artikel tidak boleh kosong", "kategori wajib dipilih", "judul artikel tidak valid untuk dijadikan slug":
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
		"message": "Artikel berhasil dibuat",
		"data":    article,
	})
}

// UpdateArticle — PUT /api/articles/:id
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
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

	article, err := h.articleSvc.UpdateArticle(id, userID, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "artikel tidak ditemukan":
			status = http.StatusNotFound
		case "akses ditolak: bukan pemilik artikel":
			status = http.StatusForbidden
		case "kategori tidak ditemukan", "satu atau lebih tag tidak ditemukan", "judul artikel tidak valid untuk dijadikan slug":
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
		"message": "Artikel berhasil diperbarui",
		"data":    article,
	})
}

// PublishArticle — POST /api/articles/:id/publish
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	article, err := h.articleSvc.PublishArticle(id)
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
		"message": "Artikel berhasil dipublish",
		"data":    article,
	})
}

// TakedownArticle — POST /api/articles/:id/takedown
func (h *ArticleHandler) TakedownArticle(c *gin.Context) {
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

	article, err := h.articleSvc.TakedownArticle(id, userID)
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
		"message": "Artikel berhasil diturunkan (takedown)",
		"data":    article,
	})
}
