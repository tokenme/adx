package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/media"
)

func mediaRouter(r *gin.Engine) {
	mediaGroup := r.Group("/media")
	mediaGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		mediaGroup.POST("/add", media.AddHandler)
		mediaGroup.GET("/verify", media.VerifyHandler)
		mediaGroup.GET("/info", media.InfoHandler)
		mediaGroup.GET("/list", media.ListHandler)
		mediaGroup.POST("/update", media.UpdateHandler)
	}
}
