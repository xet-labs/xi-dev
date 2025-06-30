package service

import (
	"xi/app/lib"
)

// xi/app/lib.* are designed so self init on method calls but adding them here ensures they are called once
func Init() {
	lib.Env.Init()
	DB.Init()
}
