// cntr/res
package cntr

import (
	"xi/app/cntr/resCntr"
	
	"github.com/gin-gonic/gin"
)

var Res = struct {
	Css gin.HandlerFunc
}{
	Css: resCntr.Css,
}

