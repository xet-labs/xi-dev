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
	Email      string     `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	ProfileImg string     `gorm:"type:varchar(255)"             json:"profile_img"`
	Role       string     `gorm:"type:varchar(50)"              json:"role"`
	Verified   bool       `json:"verified"`
	LastLogin  *time.Time `json:"last_login,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DOB        *time.Time `json:"dob,omitempty"`
	Address    string     `gorm:"type:text"                     json:"address,omitempty"`
	PhoneNo    string     `gorm:"type:varchar(20)"              json:"phone_no,omitempty"`
	Password   string     `gorm:"type:varchar(255)"             json:"-"`
	Blogs      []Blog     `gorm:"foreignKey:UID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"blogs,omitempty"` // Relation
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}
