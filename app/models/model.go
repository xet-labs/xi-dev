// models/blog
package models

import (
	"gorm.io/gorm"
	"xi/app/services"
)

// DB returns default DB or a given DB if exists
func DB(db ...string) *gorm.DB {
	return services.DB(db...)
}

