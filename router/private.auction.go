package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/privateAuction"
)

func privateAuctionRouter(r *gin.Engine) {
	privateAuctionGroup := r.Group("/private-auction")
	privateAuctionGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		privateAuctionGroup.POST("/add", privateAuction.AddHandler)
		privateAuctionGroup.POST("/update", privateAuction.UpdateHandler)
		privateAuctionGroup.GET("/list", privateAuction.ListHandler)
		privateAuctionGroup.GET("/event", privateAuction.EventHandler)
		privateAuctionGroup.POST("/audit", privateAuction.AuditHandler)
	}
}
