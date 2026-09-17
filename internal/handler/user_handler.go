package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetUsers — GET /api/users?page=1&limit=10
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.userSvc.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar user",
		"data":    result,
	})
}

// GetUser — GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	user, err := h.userSvc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data user",
		"data":    user,
	})
}

// CreateUser — POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	user, err := h.userSvc.CreateUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email sudah terdaftar" {
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
		"message": "User berhasil dibuat",
		"data":    user,
	})
}

// UpdateUser — PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	requestorRole, _ := c.Get("user_role")
	role := requestorRole.(string)

	user, err := h.userSvc.UpdateUser(id, req, role)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "user tidak ditemukan":
			status = http.StatusNotFound
		case "email sudah digunakan oleh user lain":
			status = http.StatusConflict
		case "hanya super_admin yang bisa mengangkat super_admin":
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
		"message": "User berhasil diperbarui",
		"data":    user,
	})
}

// DeleteUser — DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil {
        return
    }

    requestorID, _ := c.Get("user_id")
    myID := requestorID.(uint)

    // AMBIL ROLE DARI CONTEXT
    requestorRole, _ := c.Get("user_role")
    myRole := requestorRole.(string)

    // TAMBAHKAN myRole SEBAGAI PARAMETER KETIGA
    if err := h.userSvc.DeleteUser(id, myID, myRole); err != nil {
        status := http.StatusInternalServerError
        switch err.Error() {
        case "user tidak ditemukan":
            status = http.StatusNotFound
        case "tidak bisa menghapus akun sendiri":
            status = http.StatusBadRequest
        case "tidak punya izin menghapus akun super_admin": // Handle error role
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
        "message": "User berhasil dihapus",
    })
}

// ChangePassword — PATCH /api/users/me/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	id := userID.(uint)

	if err := h.userSvc.ChangePassword(id, req); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "password lama tidak sesuai" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password berhasil diubah",
	})
}

// AdminResetPassword — PATCH /api/users/:id/password
func (h *UserHandler) AdminResetPassword(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil {
        return
    }

    var req service.AdminChangePasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Data tidak valid",
            "errors":  err.Error(),
        })
        return
    }

    // AMBIL ROLE DARI CONTEXT
    requestorRole, _ := c.Get("user_role")
    myRole := requestorRole.(string)

    // TAMBAHKAN myRole SEBAGAI PARAMETER KETIGA
    if err := h.userSvc.AdminResetPassword(id, req, myRole); err != nil {
        status := http.StatusInternalServerError
        switch err.Error() {
        case "user tidak ditemukan":
            status = http.StatusNotFound
        case "tidak punya izin mereset password akun super_admin": // Handle error role
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
        "message": "Password user berhasil direset",
    })
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func parseUintParam(c *gin.Context, param string) (uint, error) {
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
