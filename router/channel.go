package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/channel"
)

func channelRouter(r *gin.Engine) {
	channelGroup := r.Group("/channel")
	channelGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		channelGroup.GET("/list", channel.ListGetHandler)
		channelGroup.POST("/add", channel.AddHandler)
	}
}
