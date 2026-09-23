package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create menyimpan artikel baru beserta relasi categories dan tags dalam transaksi
func (r *ArticleRepository) Create(article *model.Article, categories []model.Category, tags []model.Tag) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		article.Categories = nil // Biarkan association replace menangani categories secara eksplisit
		article.Tags = nil       // Biarkan association replace menangani tags secara eksplisit
		if err := tx.Create(article).Error; err != nil {
			return err
		}
		if len(categories) > 0 {
			if err := tx.Model(article).Association("Categories").Replace(categories); err != nil {
				return err
			}
		}
		if len(tags) > 0 {
			if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByID mencari artikel berdasarkan ID lengkap dengan Categories, Category, Tags, Author, dan TakenDownByUser
func (r *ArticleRepository) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Categories").
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Preload("Author.Role").
		Preload("TakenDownByUser").
		Preload("TakenDownByUser.Role").
		First(&article, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &article, err
}

// FindAll mengambil seluruh artikel (untuk admin / super_admin) dengan pagination
func (r *ArticleRepository) FindAll(page, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&model.Article{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Categories").
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Preload("Author.Role").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&articles).Error

	return articles, total, err
}

// FindByAuthorID mengambil daftar artikel berdasarkan author_id dengan pagination
func (r *ArticleRepository) FindByAuthorID(authorID uint, page, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&model.Article{}).Where("author_id = ?", authorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Categories").
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Preload("Author.Role").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&articles).Error

	return articles, total, err
}

// Update memperbarui data artikel dan relasi categories/tags dalam transaksi
func (r *ArticleRepository) Update(article *model.Article, categories []model.Category, updateCategories bool, tags []model.Tag, updateTags bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(article).Error; err != nil {
			return err
		}
		if updateCategories {
			if err := tx.Model(article).Association("Categories").Replace(categories); err != nil {
				return err
			}
		}
		if updateTags {
			if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		return nil
	})
}

// Save menyimpan perubahan artikel tanpa menyentuh relasi tags/categories (untuk publish/status change)
func (r *ArticleRepository) Save(article *model.Article) error {
	return r.db.Save(article).Error
}

// FindBySlug mencari artikel berdasarkan slug
func (r *ArticleRepository) FindBySlug(slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.Where("slug = ?", slug).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &article, err
}

// SlugExistsExcept mengecek apakah slug sudah digunakan artikel lain
func (r *ArticleRepository) SlugExistsExcept(slug string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Where("slug = ? AND id != ?", slug, excludeID).Count(&count).Error
	return count > 0, err
}

// FindPublishedByPostType mengambil daftar artikel published berdasarkan post_type dengan pagination
func (r *ArticleRepository) FindPublishedByPostType(postType model.PostType, page, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&model.Article{}).
		Where("status = ? AND post_type = ?", model.ArticleStatusPublished, postType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Categories").
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Preload("Author.Role").
		Offset(offset).
		Limit(limit).
		Order("published_at DESC").
		Find(&articles).Error

	return articles, total, err
}

// FindPublishedBySlugAndPostType mengambil artikel published berdasarkan slug dan post_type
func (r *ArticleRepository) FindPublishedBySlugAndPostType(slug string, postType model.PostType) (*model.Article, error) {
	var article model.Article
	err := r.db.Where("slug = ? AND status = ? AND post_type = ?", slug, model.ArticleStatusPublished, postType).
		Preload("Categories").
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Preload("Author.Role").
		First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &article, err
}

// FindPublished mengambil daftar artikel yang berstatus published dengan pagination (legacy/all)
func (r *ArticleRepository) FindPublished(page, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&model.Article{}).Where("status = ?", model.ArticleStatusPublished)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Categories").
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Preload("Author.Role").
		Offset(offset).
		Limit(limit).
		Order("published_at DESC").
		Find(&articles).Error

	return articles, total, err
}

// FindPublishedBySlug mengambil artikel published berdasarkan slug (legacy/all)
func (r *ArticleRepository) FindPublishedBySlug(slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.Where("slug = ? AND status = ?", slug, model.ArticleStatusPublished).
		Preload("Categories").
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Preload("Author.Role").
		First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &article, err
}
