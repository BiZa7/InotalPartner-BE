package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	maxImageSize = 5 * 1024 * 1024
	uploadDir    = "./uploads"
)

// UploadImage — POST /api/uploads/image
// Menerima JPG, PNG, atau WEBP dengan ukuran maksimal 5 MB.
func UploadImage(c *gin.Context) {
	// Batasi ukuran request sedikit di atas 5 MB untuk mengakomodasi multipart overhead.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxImageSize+256*1024)

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File gambar wajib dikirim dengan field 'image'",
		})
		return
	}

	if fileHeader.Size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File gambar kosong",
		})
		return
	}

	if fileHeader.Size > maxImageSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Ukuran gambar maksimal 5 MB",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuka file gambar",
		})
		return
	}
	defer file.Close()

	// Periksa MIME dari isi file, bukan dari nama/ekstensi yang dikirim client.
	var header [512]byte
	n, err := file.Read(header[:])
	if err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal membaca file gambar",
		})
		return
	}

	contentType := http.DetectContentType(header[:n])
	extension := imageExtension(contentType)
	if extension == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format gambar harus JPG, PNG, atau WEBP",
		})
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memproses file gambar",
		})
		return
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyiapkan folder upload",
		})
		return
	}

	filename, err := randomFilename(extension)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat nama file gambar",
		})
		return
	}

	destination := filepath.Join(uploadDir, filename)
	dst, err := os.Create(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyimpan gambar",
		})
		return
	}

	_, copyErr := io.Copy(dst, file)
	closeErr := dst.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(destination)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyimpan gambar",
		})
		return
	}

	url := "/uploads/" + filename

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Gambar berhasil diupload",
		"url":     url,
		"data": gin.H{
			"url": url,
		},
	})
}

func imageExtension(contentType string) string {
	switch strings.ToLower(contentType) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}

func randomFilename(extension string) (string, error) {
	var randomBytes [16]byte
	if _, err := rand.Read(randomBytes[:]); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", hex.EncodeToString(randomBytes[:]), extension), nil
}
