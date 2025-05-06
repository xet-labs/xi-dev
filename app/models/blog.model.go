// models/blog
package models

import (
	"time"

	"gorm.io/gorm"
	"xi/app/services"
)

// Blog represents the blog data model
type Blog struct {
	ID           uint       `gorm:"primaryKey;column:id"`
	UID          string     `gorm:"column:uid;index"`
	Status       string     `gorm:"column:status"`
	Tags         string     `gorm:"column:tags"`
	Title        string     `gorm:"column:title"`
	Description  string     `gorm:"column:description"`
	FeaturedImg  string     `gorm:"column:featured_img"`
	Content      string     `gorm:"column:content"`
	TitleMini    string     `gorm:"column:short_title"`
	Slug         string     `gorm:"column:slug"`
	Path         string     `gorm:"column:path"`
	MetaKeywords string     `gorm:"column:meta_keywords"`
	MetaOG       string     `gorm:"column:meta_og"`
	MetaLD       string     `gorm:"column:meta_ld"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`

	// Relation with User
	User *User `gorm:"foreignKey:UID;references:UID"`
}

// DB returns default DB or a given DB if passed
// In Blog model
func (b *Blog) DB(db ...string) *gorm.DB {
	return services.DB(db...)
}
