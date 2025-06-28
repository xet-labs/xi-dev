// cntr/res
package cntr

import (
	"xi/app/cntr/res"
	
	"github.com/gin-gonic/gin"
)

var Res = struct {
	Css gin.HandlerFunc
}{
	Css: res.Css,
}
