// models/blog
package models

import (
	"time"
)

type Blog struct {
	ID           uint       `gorm:"primaryKey;column:id"`
	UID          uint64     `gorm:"column:uid;index"`
	Status       string     `gorm:"column:status"`
	Tags         string     `gorm:"column:tags"`
	Title        string     `gorm:"column:title"`
	TitleMin     string     `gorm:"column:short_title"`
	Description  string     `gorm:"column:description"`
	Hero         string     `gorm:"column:featured_img"`
	Content      string     `gorm:"column:content"`
	Slug         string     `gorm:"column:slug"`
	Path         string     `gorm:"column:path"`
	Meta         string     `gorm:"column:meta"`
	MetaKeywords string     `gorm:"column:meta_keywords"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`

	// Relation with User
	User *User `gorm:"foreignKey:UID;references:UID"`
}
