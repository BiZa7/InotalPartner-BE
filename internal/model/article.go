package model

import (
	"time"
)

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
	ArticleStatusTakenDown ArticleStatus = "taken_down"
)

type PostType string

const (
	PostTypeNews    PostType = "news"
	PostTypeEvent   PostType = "event"
	PostTypeArticle PostType = "article"
)

type Article struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `                                json:"created_at"`
	UpdatedAt time.Time `                                json:"updated_at"`

	// Post Type
	PostType PostType `gorm:"type:varchar(20);not null;default:'news'" json:"post_type"`

	// Content
	Title        string `gorm:"type:varchar(255);not null"         json:"title"`
	Slug         string `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content      string `gorm:"type:text;not null"                 json:"content"`
	Excerpt      string `gorm:"type:text"                          json:"excerpt,omitempty"`
	ThumbnailURL string `gorm:"type:text"                          json:"thumbnail_url,omitempty"`

	// Event Details (Nullable, only used for post_type = event)
	EventDate     *string `gorm:"type:date" json:"event_date,omitempty"`
	EventTime     *string `gorm:"type:time" json:"event_time,omitempty"`
	EventLocation *string `gorm:"type:text" json:"event_location,omitempty"`

	// Category Relation (1 Category -> N Articles, legacy)
	CategoryID uint     `gorm:"index"                                                             json:"category_id,omitempty"`
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"category,omitempty"`

	// Multiple Categories Relation (Many-to-Many via article_categories)
	Categories []Category `gorm:"many2many:article_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"categories,omitempty"`

	// Author Relation (1 User -> N Articles)
	AuthorID uint `gorm:"not null;index"                                                    json:"author_id"`
	Author   User `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"author,omitempty"`

	// Status & Lifecycle
	Status          ArticleStatus `gorm:"type:varchar(20);not null;default:'draft'"                        json:"status"`
	PublishedAt     *time.Time    `                                                                         json:"published_at,omitempty"`
	TakenDownAt     *time.Time    `                                                                         json:"taken_down_at,omitempty"`
	TakenDownBy     *uint         `gorm:"index"                                                             json:"taken_down_by,omitempty"`
	TakenDownByUser *User         `gorm:"foreignKey:TakenDownBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"taken_down_by_user,omitempty"`

	// Tags Relation (Many-to-Many via article_tags)
	Tags []Tag `gorm:"many2many:article_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"tags,omitempty"`
}

