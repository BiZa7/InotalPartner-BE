package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole mengizinkan akses hanya untuk role yang ditentukan.
// Gunakan setelah AuthMiddleware.
//
// Contoh:
//
//	protected.Use(middleware.RequireRole("super_admin", "admin"))
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Autentikasi diperlukan",
			})
			return
		}

		if _, ok := allowed[role.(string)]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak: hak akses tidak mencukupi",
			})
			return
		}

		c.Next()
	}
}

// SuperAdminOnly alias RequireRole("super_admin")
func SuperAdminOnly() gin.HandlerFunc {
	return RequireRole("super_admin")
}

// AdminOrAbove mengizinkan super_admin dan admin
func AdminOrAbove() gin.HandlerFunc {
	return RequireRole("super_admin", "admin")
}

// OperatorOnly mengizinkan hanya operator
func OperatorOnly() gin.HandlerFunc {
	return RequireRole("operator")
}

// DenyAll melarang akses untuk semua role
func DenyAll() gin.HandlerFunc {
	return RequireRole()
}
