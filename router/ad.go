package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/ad"
)

func adRouter(r *gin.Engine) {
	r.GET("/c", ad.ShowHandler)
	r.GET("/t/:key", ad.ClickHandler)
	r.GET("/i/:key", ad.ImpHandler)
	r.GET("/ad/loader", ad.LoaderHandler)
}
