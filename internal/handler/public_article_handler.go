package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type PublicArticleHandler struct {
	articleSvc *service.ArticleService
}

func NewPublicArticleHandler(articleSvc *service.ArticleService) *PublicArticleHandler {
	return &PublicArticleHandler{articleSvc: articleSvc}
}

// GetPublicArticles — GET /api/public/articles?page=1&limit=10
// Public endpoint — tidak membutuhkan autentikasi
func (h *PublicArticleHandler) GetPublicArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.articleSvc.GetPublishedArticles(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar artikel publik",
		"data":    result,
	})
}

// GetPublicArticleBySlug — GET /api/public/articles/:slug
// Public endpoint — tidak membutuhkan autentikasi
func (h *PublicArticleHandler) GetPublicArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Slug artikel tidak valid",
		})
		return
	}

	article, err := h.articleSvc.GetPublishedArticleBySlug(slug)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "artikel tidak ditemukan" {
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
		"message": "Berhasil mengambil data artikel",
		"data":    article,
	})
}
