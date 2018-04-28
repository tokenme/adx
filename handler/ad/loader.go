package ad

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func LoaderHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, fmt.Sprintf("%stmm.%s.js", Config.CDNUrl, Config.AdJSVer))
	return
}
