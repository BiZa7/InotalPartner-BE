package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/model"
	"inotal-be/internal/service"
)

type EventHandler struct {
	articleSvc *service.ArticleService
}

func NewEventHandler(articleSvc *service.ArticleService) *EventHandler {
	return &EventHandler{articleSvc: articleSvc}
}

// GetEvents — GET /api/events?page=1&limit=10
func (h *EventHandler) GetEvents(c *gin.Context) {
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

	result, err := h.articleSvc.GetArticlesByPostType(userID, userRole, model.PostTypeEvent, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar event",
		"data":    result,
	})
}

// GetEventByID — GET /api/events/:id
func (h *EventHandler) GetEventByID(c *gin.Context) {
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

	article, err := h.articleSvc.GetArticleByIDAndPostType(id, userID, userRole, model.PostTypeEvent)
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
		"message": "Berhasil mengambil data event",
		"data":    article,
	})
}

// CreateEvent — POST /api/events
func (h *EventHandler) CreateEvent(c *gin.Context) {
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

	article, err := h.articleSvc.CreateArticleWithPostType(userID, req, model.PostTypeEvent)
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
		"message": "Event berhasil dibuat",
		"data":    article,
	})
}

// UpdateEvent — PUT /api/events/:id
func (h *EventHandler) UpdateEvent(c *gin.Context) {
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

	article, err := h.articleSvc.UpdateArticleWithPostType(id, userID, req, model.PostTypeEvent)
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
		"message": "Event berhasil diperbarui",
		"data":    article,
	})
}

// PublishEvent — POST /api/events/:id/publish
func (h *EventHandler) PublishEvent(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	article, err := h.articleSvc.PublishArticleWithPostType(id, model.PostTypeEvent)
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
		"message": "Event berhasil dipublish",
		"data":    article,
	})
}

// TakedownEvent — POST /api/events/:id/takedown
func (h *EventHandler) TakedownEvent(c *gin.Context) {
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

	article, err := h.articleSvc.TakedownArticleWithPostType(id, userID, model.PostTypeEvent)
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
		"message": "Event berhasil diturunkan (takedown)",
		"data":    article,
	})
}
