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

        // Ambil role TERKINI dari database, bukan dari token.
        // Role di dalam token adalah "foto" saat login; kalau role user diubah
        // (mis. dipromosikan jadi admin) setelah itu, token lama tetap membawa
        // role basi sampai user login ulang. Dengan mengambil ulang dari DB di
        // sini, perubahan role langsung berlaku di request berikutnya.
        currentRole, err := authSvc.CurrentRole(claims.UserID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "message": "Gagal memverifikasi akun",
            })
            return
        }
        if currentRole == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "Akun tidak ditemukan atau sudah dinonaktifkan",
            })
            return
        }

        // Simpan claims ke context agar handler bisa mengakses
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_role", currentRole)
        c.Next()
    }
}