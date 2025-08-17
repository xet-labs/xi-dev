// models/blog
package model

import (
	"time"
	"xi/app/model/util"
)

type Blog struct {
	ID           uint             `gorm:"primaryKey;column:id"          json:"id"`
	UID          uint             `gorm:"column:uid;index"              json:"uid"`
	Status       string           `gorm:"column:status"                 json:"status"`
	Tags         util.StringArray `gorm:"column:tags"                   json:"tags"`
	Title        string           `gorm:"column:title"                  json:"title"`
	Headline     string           `gorm:"column:short_title"            json:"headline"`
	Description  string           `gorm:"column:description"            json:"description"`
	FeaturedImg  string           `gorm:"column:featured_img"           json:"featured_img"`
	Content      string           `gorm:"column:content"                json:"content"`
	Slug         string           `gorm:"column:slug"                   json:"slug"`
	Path         string           `gorm:"column:path"                   json:"path"`
	Meta         string           `gorm:"column:meta"                   json:"meta"`
	MetaKeywords string           `gorm:"column:meta_keywords"          json:"meta_keywords"`
	CreatedAt    *time.Time       `gorm:"column:created_at"             json:"created_at"`
	UpdatedAt    *time.Time       `gorm:"column:updated_at"             json:"updated_at"`
	User         *User            `gorm:"foreignKey:UID;references:UID" json:"user,omitempty"` // Relation
}
