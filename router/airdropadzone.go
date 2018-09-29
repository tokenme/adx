package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/airdropadzone"
)

func airdropAdzoneRouter(r *gin.Engine) {
	g := r.Group("/airdropadzone")
	g.Use(AuthMiddleware.MiddlewareFunc())
	{
		g.GET("/list", airdropadzone.ListGetHandler)
		g.POST("/add", airdropadzone.AddHandler)
	}
}
