package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/adzone"
)

func adzoneRouter(r *gin.Engine) {
	adzoneGroup := r.Group("/adzone")
	adzoneGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		adzoneGroup.POST("/add", adzone.AddHandler)
		adzoneGroup.GET("/list", adzone.ListHandler)
		adzoneGroup.GET("/info", adzone.InfoHandler)
		adzoneGroup.GET("/search", adzone.SearchHandler)
		adzoneGroup.POST("/update", adzone.UpdateHandler)
		adzoneGroup.GET("/Trafficlist",adzone.TrafficListHandler)
		adzoneGroup.GET("/MediaList",adzone.MediaListHandler)
		//adzoneGroup.GET("/info", adzone.InfoHandler)
	}
	r.GET("/adzone/sizes", adzone.SizeListHandler)
}
