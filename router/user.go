package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/user"
)

func userRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		userGroup.GET("/info", user.InfoGetHandler)
		userGroup.GET("/balance", user.BalanceHandler)
	}
	r.GET("/user/activate", user.ActivateHandler)
	r.GET("/user/reset-pwd-verify", user.ResetPwdVerifyHandler)
	r.POST("/user/create", user.CreateHandler)
	r.POST("/user/reset-password", user.ResetPasswordHandler)
	r.GET("/user/avatar/:key", user.AvatarGetHandler)
	r.POST("/user/resend-activation-email", user.ResendActivationEmailHandler)
	r.POST("/user/send-reset-password-email", user.SendResetPasswordEmailHandler)
}
