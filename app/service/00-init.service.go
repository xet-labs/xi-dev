// services/env
package service

import (
	"xi/app/util"
)

func Init() {
	util.InitEnv()
	InitDB()
}
