package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/stats"
)

func statsRouter(r *gin.Engine) {
	statsGroup := r.Group("/stats")
	statsGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		statsGroup.GET("/dates", stats.DatesHandler)
		statsGroup.GET("/country", stats.CountryHandler)
		statsGroup.GET("/browser-type", stats.BrowserTypeHandler)
		statsGroup.GET("/os", stats.OsHandler)
		statsGroup.GET("/browser", stats.BrowserHandler)
	}
}
