package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/airdrop"
)

func airdropRouter(r *gin.Engine) {
	airdropGroup := r.Group("/airdrop")
	airdropGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		airdropGroup.POST("/add", airdrop.AddHandler)
		airdropGroup.GET("/info", airdrop.InfoHandler)
		airdropGroup.GET("/list", airdrop.ListHandler)
		airdropGroup.POST("/update", airdrop.UpdateHandler)
	}
	airdropCheckGroup := r.Group("/airdrop")
	airdropCheckGroup.Use(AuthCheckerFunc())
	airdropCheckGroup.POST("/share", airdrop.ShareHandler)
}
