package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateArticleRequest struct {
	Title        string `json:"title" binding:"required,min=1,max=255"`
	Content      string `json:"content" binding:"required"`
	CategoryID   uint   `json:"category_id" binding:"required"`
	TagIDs       []uint `json:"tag_ids"`
	Excerpt      string `json:"excerpt"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type UpdateArticleRequest struct {
	Title        string  `json:"title" binding:"omitempty,min=1,max=255"`
	Content      string  `json:"content"`
	CategoryID   uint    `json:"category_id"`
	TagIDs       *[]uint `json:"tag_ids"`
	Excerpt      string  `json:"excerpt"`
	ThumbnailURL string  `json:"thumbnail_url"`
}

type ArticleListResponse struct {
	Data       []model.Article `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
}

func NewArticleService(
	articleRepo *repository.ArticleRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// GetArticles mengambil daftar artikel (Operator: milik sendiri, Admin/SuperAdmin: semua)
func (s *ArticleService) GetArticles(userID uint, userRole string, page, limit int) (*ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var articles []model.Article
	var total int64
	var err error

	if userRole == "admin" || userRole == "super_admin" {
		articles, total, err = s.articleRepo.FindAll(page, limit)
	} else {
		articles, total, err = s.articleRepo.FindByAuthorID(userID, page, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &ArticleListResponse{
		Data:       articles,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetArticleByID mengambil detail artikel (Operator: milik sendiri, Admin/SuperAdmin: semua)
func (s *ArticleService) GetArticleByID(id uint, userID uint, userRole string) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if article == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	// Ownership check untuk Operator
	if userRole != "admin" && userRole != "super_admin" {
		if article.AuthorID != userID {
			return nil, errors.New("akses ditolak: bukan pemilik artikel")
		}
	}

	return article, nil
}

// CreateArticle membuat artikel baru milik operator dengan status draft
func (s *ArticleService) CreateArticle(authorID uint, req CreateArticleRequest) (*model.Article, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("judul artikel tidak boleh kosong")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("konten artikel tidak boleh kosong")
	}

	if req.CategoryID == 0 {
		return nil, errors.New("kategori wajib dipilih")
	}

	// Validasi category exists
	category, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("database error saat validasi kategori: %w", err)
	}
	if category == nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	// Validasi tag_ids
	var tags []model.Tag
	if len(req.TagIDs) > 0 {
		uniqueTagIDs := make([]uint, 0, len(req.TagIDs))
		seen := make(map[uint]bool)
		for _, tid := range req.TagIDs {
			if !seen[tid] {
				seen[tid] = true
				uniqueTagIDs = append(uniqueTagIDs, tid)
			}
		}

		foundTags, err := s.tagRepo.FindByIDs(uniqueTagIDs)
		if err != nil {
			return nil, fmt.Errorf("database error saat validasi tag: %w", err)
		}
		if len(foundTags) != len(uniqueTagIDs) {
			return nil, errors.New("satu atau lebih tag tidak ditemukan")
		}
		tags = foundTags
	}

	// Generate unique slug
	baseSlug := generateSlug(title)
	if baseSlug == "" {
		return nil, errors.New("judul artikel tidak valid untuk dijadikan slug")
	}

	slug, err := generateUniqueSlug(baseSlug, func(candidate string) (bool, error) {
		found, err := s.articleRepo.FindBySlug(candidate)
		return found != nil, err
	})
	if err != nil {
		return nil, err
	}

	article := &model.Article{
		Title:        title,
		Slug:         slug,
		Content:      req.Content,
		Excerpt:      strings.TrimSpace(req.Excerpt),
		ThumbnailURL: strings.TrimSpace(req.ThumbnailURL),
		CategoryID:   req.CategoryID,
		AuthorID:     authorID,
		Status:       model.ArticleStatusDraft,
	}

	if err := s.articleRepo.Create(article, tags); err != nil {
		return nil, fmt.Errorf("gagal membuat artikel: %w", err)
	}

	// Ambil data artikel lengkap dengan relasi untuk dikembalikan
	return s.articleRepo.FindByID(article.ID)
}

// UpdateArticle memperbarui artikel (Hanya Operator yang boleh memperbarui artikel miliknya sendiri)
func (s *ArticleService) UpdateArticle(id uint, authorID uint, req UpdateArticleRequest) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if article == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	// Ownership check (Hanya pemilik artikel)
	if article.AuthorID != authorID {
		return nil, errors.New("akses ditolak: bukan pemilik artikel")
	}

	// Validasi dan update Title & Slug
	if req.Title != "" {
		newTitle := strings.TrimSpace(req.Title)
		if newTitle != "" && newTitle != article.Title {
			baseSlug := generateSlug(newTitle)
			if baseSlug == "" {
				return nil, errors.New("judul artikel tidak valid untuk dijadikan slug")
			}

			slug, err := generateUniqueSlug(baseSlug, func(candidate string) (bool, error) {
				exists, err := s.articleRepo.SlugExistsExcept(candidate, article.ID)
				return exists, err
			})
			if err != nil {
				return nil, err
			}
			article.Title = newTitle
			article.Slug = slug
		}
	}

	// Update Content jika dikirim
	if req.Content != "" {
		article.Content = req.Content
	}

	// Update Excerpt jika dikirim
	if req.Excerpt != "" {
		article.Excerpt = strings.TrimSpace(req.Excerpt)
	}

	// Update ThumbnailURL jika dikirim
	if req.ThumbnailURL != "" {
		article.ThumbnailURL = strings.TrimSpace(req.ThumbnailURL)
	}

	// Validasi dan update Category jika dikirim
	if req.CategoryID != 0 && req.CategoryID != article.CategoryID {
		category, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("database error saat validasi kategori: %w", err)
		}
		if category == nil {
			return nil, errors.New("kategori tidak ditemukan")
		}
		article.CategoryID = req.CategoryID
	}

	// Validasi dan update Tags jika dikirim
	var tags []model.Tag
	updateTags := req.TagIDs != nil
	if updateTags {
		tagIDs := *req.TagIDs
		if len(tagIDs) > 0 {
			uniqueTagIDs := make([]uint, 0, len(tagIDs))
			seen := make(map[uint]bool)
			for _, tid := range tagIDs {
				if !seen[tid] {
					seen[tid] = true
					uniqueTagIDs = append(uniqueTagIDs, tid)
				}
			}

			foundTags, err := s.tagRepo.FindByIDs(uniqueTagIDs)
			if err != nil {
				return nil, fmt.Errorf("database error saat validasi tag: %w", err)
			}
			if len(foundTags) != len(uniqueTagIDs) {
				return nil, errors.New("satu atau lebih tag tidak ditemukan")
			}
			tags = foundTags
		}
	}

	// Note: Status dan AuthorID tidak diubah di sini
	if err := s.articleRepo.Update(article, tags, updateTags); err != nil {
		return nil, fmt.Errorf("gagal memperbarui artikel: %w", err)
	}

	return s.articleRepo.FindByID(id)
}

// PublishArticle mem-publish artikel draft (hanya Admin / Super Admin)
func (s *ArticleService) PublishArticle(id uint) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if article == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	// Status check
	switch article.Status {
	case model.ArticleStatusPublished:
		return nil, errors.New("artikel sudah dipublish")
	case model.ArticleStatusTakenDown:
		return nil, errors.New("artikel yang sudah diturunkan tidak dapat dipublish kembali")
	}

	// Set publish
	now := time.Now()
	article.Status = model.ArticleStatusPublished
	article.PublishedAt = &now

	if err := s.articleRepo.Save(article); err != nil {
		return nil, fmt.Errorf("gagal mempublish artikel: %w", err)
	}

	return s.articleRepo.FindByID(id)
}

// TakedownArticle men-takedown artikel yang berstatus published (hanya Admin / Super Admin)
func (s *ArticleService) TakedownArticle(id uint, adminUserID uint) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if article == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if article.Status != model.ArticleStatusPublished {
		if article.Status == model.ArticleStatusTakenDown {
			return nil, errors.New("artikel sudah diturunkan")
		}
		return nil, errors.New("hanya artikel yang berstatus published yang dapat diturunkan")
	}

	now := time.Now()
	article.Status = model.ArticleStatusTakenDown
	article.TakenDownAt = &now
	article.TakenDownBy = &adminUserID

	if err := s.articleRepo.Save(article); err != nil {
		return nil, fmt.Errorf("gagal men-takedown artikel: %w", err)
	}

	return s.articleRepo.FindByID(id)
}

// ─── Public (No Auth) ─────────────────────────────────────────────────────────

// GetPublishedArticles mengambil daftar artikel published untuk publik dengan pagination
func (s *ArticleService) GetPublishedArticles(page, limit int) (*ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	articles, total, err := s.articleRepo.FindPublished(page, limit)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &ArticleListResponse{
		Data:       articles,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetPublishedArticleBySlug mengambil detail artikel published berdasarkan slug untuk publik
func (s *ArticleService) GetPublishedArticleBySlug(slug string) (*model.Article, error) {
	article, err := s.articleRepo.FindPublishedBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if article == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}
	return article, nil
}
