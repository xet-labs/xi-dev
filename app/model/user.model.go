// models/user
package model

import (
	"time"
)

type User struct {
	UID        uint       `gorm:"primaryKey"                    json:"uid"`
	Status     string     `gorm:"type:varchar(50)"              json:"status"`
	Username   string     `gorm:"type:varchar(255);uniqueIndex" json:"username"`
	Name       string     `gorm:"type:varchar(255)"             json:"name"`
	Email      string     `gorm:"type:varchar(255);uniqueIndex" json:"email"                 binding:"email"`
	ProfileImg string     `gorm:"type:varchar(255)"             json:"profile_img"`
	Role       string     `gorm:"type:varchar(50)"              json:"role"`
	Verified   bool       `json:"verified"`
	LastLogin  *time.Time `gorm:"type:timestamp"                json:"last_login,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at"             json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"             json:"updated_at"`
	DOB        *time.Time `gorm:"type:date"                     json:"dob,omitempty"`
	Address    string     `gorm:"type:text"                     json:"address,omitempty"`
	PhoneNo    string     `gorm:"type:varchar(20)"              json:"phone_no,omitempty"`
	Password   string     `gorm:"type:varchar(255)"             json:"-"`
	Blogs      []Blog     `gorm:"foreignKey:UID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"blogs,omitempty"` // Relation
}
