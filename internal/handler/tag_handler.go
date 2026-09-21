package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type TagHandler struct {
	tagSvc *service.TagService
}

func NewTagHandler(tagSvc *service.TagService) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}

// GetTags — GET /api/tags
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.tagSvc.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar tag",
		"data":    tags,
	})
}

// GetTag — GET /api/tags/:id
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	tag, err := h.tagSvc.GetTagByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "tag tidak ditemukan" {
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
		"message": "Berhasil mengambil data tag",
		"data":    tag,
	})
}

// CreateTag — POST /api/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req service.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	tag, err := h.tagSvc.CreateTag(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "tag sudah ada" {
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
		"message": "Tag berhasil dibuat",
		"data":    tag,
	})
}

// UpdateTag — PUT /api/tags/:id
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	tag, err := h.tagSvc.UpdateTag(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "tag tidak ditemukan":
			status = http.StatusNotFound
		case "tag sudah ada":
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
		"message": "Tag berhasil diperbarui",
		"data":    tag,
	})
}

// DeleteTag — DELETE /api/tags/:id
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.tagSvc.DeleteTag(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "tag tidak ditemukan" {
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
		"message": "Tag berhasil dihapus",
	})
}
