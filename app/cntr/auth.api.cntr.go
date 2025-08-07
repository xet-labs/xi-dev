// cntr/blog.api.go
package cntr

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthApiCntr struct {
	db    *gorm.DB
	rdb   *redis.Client
}

// Singleton controller
var AuthApi = &BlogApiCntr{}


