package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/token"
)

func tokenRouter(r *gin.Engine) {
	tokenGroup := r.Group("/token")
	tokenGroup.GET("/get", token.GetHandler)
}
