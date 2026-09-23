package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/model"
	"inotal-be/internal/service"
)

type PublicArticleHandler struct {
	articleSvc *service.ArticleService
}

func NewPublicArticleHandler(articleSvc *service.ArticleService) *PublicArticleHandler {
	return &PublicArticleHandler{articleSvc: articleSvc}
}

// ── NEWS ─────────────────────────────────────────────────────────────────────

// GetPublicNews — GET /api/public/news?page=1&limit=10
func (h *PublicArticleHandler) GetPublicNews(c *gin.Context) {
	h.getPublicList(c, model.PostTypeNews, "Berhasil mengambil daftar berita publik")
}

// GetPublicNewsBySlug — GET /api/public/news/:slug
func (h *PublicArticleHandler) GetPublicNewsBySlug(c *gin.Context) {
	h.getPublicBySlug(c, model.PostTypeNews, "Berhasil mengambil data berita")
}

// ── EVENTS ───────────────────────────────────────────────────────────────────

// GetPublicEvents — GET /api/public/events?page=1&limit=10
func (h *PublicArticleHandler) GetPublicEvents(c *gin.Context) {
	h.getPublicList(c, model.PostTypeEvent, "Berhasil mengambil daftar event publik")
}

// GetPublicEventsBySlug — GET /api/public/events/:slug
func (h *PublicArticleHandler) GetPublicEventsBySlug(c *gin.Context) {
	h.getPublicBySlug(c, model.PostTypeEvent, "Berhasil mengambil data event")
}

// ── ARTICLES ─────────────────────────────────────────────────────────────────

// GetPublicArticles — GET /api/public/articles?page=1&limit=10
func (h *PublicArticleHandler) GetPublicArticles(c *gin.Context) {
	h.getPublicList(c, model.PostTypeArticle, "Berhasil mengambil daftar artikel publik")
}

// GetPublicArticleBySlug — GET /api/public/articles/:slug
func (h *PublicArticleHandler) GetPublicArticleBySlug(c *gin.Context) {
	h.getPublicBySlug(c, model.PostTypeArticle, "Berhasil mengambil data artikel")
}

// ── HELPERS ──────────────────────────────────────────────────────────────────

func (h *PublicArticleHandler) getPublicList(c *gin.Context, postType model.PostType, successMsg string) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.articleSvc.GetPublishedArticlesByPostType(postType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": successMsg,
		"data":    result,
	})
}

func (h *PublicArticleHandler) getPublicBySlug(c *gin.Context, postType model.PostType, successMsg string) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Slug tidak valid",
		})
		return
	}

	article, err := h.articleSvc.GetPublishedArticleBySlugAndPostType(slug, postType)
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
		"message": successMsg,
		"data":    article,
	})
}
