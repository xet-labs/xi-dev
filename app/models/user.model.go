// models/user
package models

import (
	"time"

	"gorm.io/gorm"
	"xi/app/services"
	"xi/app/utils"
)

// User represents a user in the system
type User struct {
	UID        uint64   `gorm:"primaryKey"`
	Status     string `gorm:"type:varchar(50)"`
	Username   string `gorm:"type:varchar(255);uniqueIndex"`
	Name       string `gorm:"type:varchar(255)"`
	Email      string `gorm:"type:varchar(255);uniqueIndex"`
	ProfileImg string `gorm:"type:varchar(255)"`
	Role       string `gorm:"type:varchar(50)"`
	Verified   bool
	LastLogin  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DOB        *time.Time
	Address    string `gorm:"type:text"`
	PhoneNo    string `gorm:"type:varchar(20)"`
	Password   string `gorm:"type:varchar(255)" json:"-"` // Hide password from JSON

	// Relationships
	Blogs []Blog `gorm:"foreignKey:UID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	
}

// DB returns default DB or a given DB if passed
func (u *User) DB(db ...*gorm.DB) *gorm.DB {
	if len(db) > 0 && db[0] != nil {
		return db[0]
	}
	defaultDB := utils.Env("DB_DEFAULT", "XI")

	return services.DB(defaultDB)
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}
