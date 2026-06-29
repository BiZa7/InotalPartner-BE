package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

// AuthMiddleware memvalidasi JWT dari header Authorization: Bearer <token>
func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header diperlukan",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Format token tidak valid, gunakan: Bearer <token>",
			})
			return
		}

		claims, err := authSvc.ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak valid atau sudah expired",
			})
			return
		}

		// Simpan claims ke context agar handler bisa mengakses
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", string(claims.Role))
		c.Next()
	}
}

// AdminOnly hanya izinkan role admin — gunakan setelah AuthMiddleware
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak: hanya admin",
			})
			return
		}
		c.Next()
	}
}