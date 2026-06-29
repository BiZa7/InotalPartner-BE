package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/config"
	"inotal-be/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Register
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	resp, err := h.authSvc.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data":    resp,
	})
}

// Login
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	resp, err := h.authSvc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    resp,
	})
}

// Google OAuth: Redirect
// GET /api/auth/google
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Generate random state untuk CSRF protection
	state := generateState()
	// Simpan state di cookie sementara (30 menit)
	c.SetCookie("oauth_state", state, 1800, "/", "", false, true)

	url := h.authSvc.GoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Google OAuth: Callback
// GET /api/auth/google/callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Validasi state (CSRF check)
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || cookieState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid OAuth state",
		})
		return
	}
	// Hapus cookie state
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Authorization code tidak ditemukan",
		})
		return
	}

	resp, err := h.authSvc.GoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Google login gagal: " + err.Error(),
		})
		return
	}

	// Redirect ke frontend dengan token di query param
	// Frontend menyimpan token ke localStorage
	frontendURL := config.App.FrontendURL
	c.Redirect(http.StatusTemporaryRedirect,
		frontendURL+"/auth/callback?token="+resp.Token,
	)
}

// Me (Protected) 
// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userEmail, _ := c.Get("user_email")
	userRole, _ := c.Get("user_role")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":    userID,
			"email": userEmail,
			"role":  userRole,
		},
	})
}

// Helper 
func generateState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}