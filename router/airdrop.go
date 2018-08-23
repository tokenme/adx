package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/airdrop"
	"github.com/tokenme/adx/middlewares/gzip"
)

func airdropRouter(r *gin.Engine) {
	airdropGroup := r.Group("/airdrop")
	airdropGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		airdropGroup.POST("/add", airdrop.AddHandler)
		airdropGroup.GET("/list", airdrop.ListHandler)
		airdropGroup.POST("/update", airdrop.UpdateHandler)

		airdropGroup.GET("/get", airdrop.GetHandler)
		airdropGroup.GET("/stats", airdrop.StatsHandler)
		airdropGroup.GET("/publisher/apply", airdrop.PublisherApplyHandler)
		airdropGroup.POST("/withdraw", airdrop.WithdrawHandler)
	}
	//airdropCheckGroup := r.Group("/airdrop")
	//airdropCheckGroup.Use(AuthCheckerFunc())
	//airdropCheckGroup.POST("/share", airdrop.ShareHandler)
	r.GET("/airdrop/submission-export", AuthMiddleware.MiddlewareFunc(), gzip.Gzip(gzip.BestCompression), airdrop.SubmissionExportHandler)
}
