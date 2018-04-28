package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/creative"
)

func creativeRouter(r *gin.Engine) {
	creativeGroup := r.Group("/creative")
	creativeGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		creativeGroup.POST("/upload", creative.UploadHandler)
	}
}
