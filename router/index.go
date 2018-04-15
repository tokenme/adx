package router

import (
	"github.com/danielkov/gin-helmet"
	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/router/static"
)

func NewRouter(uiPath string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(helmet.Default())
	xssMdlwr := &xss.XssMw{
		FieldsToSkip: []string{"password", "start_date", "end_date", "token"},
		BmPolicy:     "UGCPolicy",
	}
	r.Use(xssMdlwr.RemoveXss())
	r.Use(static.Serve("/", static.LocalFile(uiPath, 0, true)))
	r.LoadHTMLGlob("templates/*")
	authRouter(r)
	userRouter(r)
	mediaRouter(r)
	adzoneRouter(r)
	return r
}