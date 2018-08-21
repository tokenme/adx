package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/helper"
)

func geoIPRouter(r *gin.Engine) {

	r.GET("/geoip", helper.GeoIPHandler)
}
