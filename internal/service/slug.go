package service

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// spaceOrPunct digunakan untuk mengganti spasi dan karakter non-alfanumerik menjadi "-"
var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

// generateSlug membuat slug dari string input.
// Contoh: "Teknologi Informasi" → "teknologi-informasi"
// Contoh: "Kesehatan & Masyarakat" → "kesehatan-masyarakat"
func generateSlug(input string) string {
	// Normalize unicode (NFD) lalu buang karakter combining (aksen, dsb.)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, _ := transform.String(t, input)

	// Lowercase
	slug := strings.ToLower(normalized)

	// Ganti semua karakter non-alphanumeric (termasuk spasi, simbol) dengan "-"
	slug = nonAlphanumericRegex.ReplaceAllString(slug, "-")

	// Trim "-" di awal dan akhir
	slug = strings.Trim(slug, "-")

	return slug
}

// generateUniqueSlug membuat slug yang unique berdasarkan base slug.
// Jika base slug sudah ada di database (via existsFn), tambahkan suffix angka.
// Contoh: "teknologi" sudah ada → coba "teknologi-2", "teknologi-3", dst.
func generateUniqueSlug(base string, existsFn func(slug string) (bool, error)) (string, error) {
	candidate := base
	for i := 2; i <= 100; i++ {
		exists, err := existsFn(candidate)
		if err != nil {
			return "", fmt.Errorf("database error saat cek slug: %w", err)
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return "", fmt.Errorf("tidak dapat membuat slug unique untuk: %s", base)
}
