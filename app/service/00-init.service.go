// services/env
package service

import (
	"xi/app/lib"
)

func Init() {
	lib.InitEnv()
	InitDB()
}
