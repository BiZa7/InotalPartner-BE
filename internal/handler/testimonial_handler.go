package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type TestimonialHandler struct {
	testimonialSvc *service.TestimonialService
}

func NewTestimonialHandler(testimonialSvc *service.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{
		testimonialSvc: testimonialSvc,
	}
}

// ─── API Publik ───────────────────────────────────────────────────────────────

// GetPublicTestimonials — GET /api/public/testimonials (Public, No Auth)
// Menampilkan seluruh testimoni yang aktif, diurutkan berdasarkan display_order ASC, id ASC.
func (h *TestimonialHandler) GetPublicTestimonials(c *gin.Context) {
	testimonials, err := h.testimonialSvc.GetActiveTestimonials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar testimoni",
		"data":    testimonials,
	})
}

// ─── API Pengelolaan Admin ───────────────────────────────────────────────────

// GetTestimonials — GET /api/testimonials (Admin Only)
// Menampilkan seluruh testimoni, termasuk yang nonaktif.
func (h *TestimonialHandler) GetTestimonials(c *gin.Context) {
	testimonials, err := h.testimonialSvc.GetAllTestimonials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh daftar testimoni",
		"data":    testimonials,
	})
}

// GetTestimonialByID — GET /api/testimonials/:id (Admin Only)
// Menampilkan detail testimoni berdasarkan ID.
func (h *TestimonialHandler) GetTestimonialByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	testimonial, err := h.testimonialSvc.GetTestimonialByID(id)
	if err != nil {
		if errors.Is(err, service.ErrTestimonialNotFound) {
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
		"message": "Berhasil mengambil detail testimoni",
		"data":    testimonial,
	})
}

// CreateTestimonial — POST /api/testimonials (Admin Only)
// Membuat record testimoni baru.
func (h *TestimonialHandler) CreateTestimonial(c *gin.Context) {
	var req service.CreateTestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	testimonial, err := h.testimonialSvc.CreateTestimonial(req)
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
		"message": "Testimoni berhasil ditambahkan",
		"data":    testimonial,
	})
}

// UpdateTestimonial — PUT /api/testimonials/:id (Admin Only)
// Memperbarui data testimoni berdasarkan ID.
func (h *TestimonialHandler) UpdateTestimonial(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateTestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	testimonial, err := h.testimonialSvc.UpdateTestimonial(id, req)
	if err != nil {
		if errors.Is(err, service.ErrTestimonialNotFound) {
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
		"message": "Testimoni berhasil diperbarui",
		"data":    testimonial,
	})
}

// DeleteTestimonial — DELETE /api/testimonials/:id (Admin Only)
// Menghapus record testimoni berdasarkan ID.
func (h *TestimonialHandler) DeleteTestimonial(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.testimonialSvc.DeleteTestimonial(id); err != nil {
		if errors.Is(err, service.ErrTestimonialNotFound) {
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
		"message": "Testimoni berhasil dihapus",
	})
}
