package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/auth"
	"github.com/tokenme/adx/middlewares/jwt"
	"time"
)

var AUTH_KEY = []byte("20eefe8d82ba3ca8a417e14a48d24632bc35bbd7")

const (
	AUTH_REALM      = "Tokenme.Server[tokenme.io]"
	AUTH_TIMEOUT    = 168 * time.Hour
	AUTH_MAXREFRESH = 1 * time.Hour
)

var AuthMiddleware = &jwt.GinJWTMiddleware{
	Realm:         AUTH_REALM,
	Key:           AUTH_KEY,
	Timeout:       AUTH_TIMEOUT,
	MaxRefresh:    AUTH_MAXREFRESH,
	Authenticator: auth.AuthenticatorFunc,
	Authorizator:  auth.AuthorizatorFunc,
	Unauthorized: func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    401,
			"message": message,
		})
	},
	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup: "header:Authorization",
	// TokenLookup: "query:token",
	// TokenLookup: "cookie:token",

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName: "Bearer",

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc: time.Now,
}

func AuthCheckerFunc() gin.HandlerFunc {
	if err := AuthMiddleware.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			c.Next()
			return
		}
	}

	return func(c *gin.Context) {
		token, err := AuthMiddleware.ParseToken(c)

		if err != nil {
			c.Next()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		id := AuthMiddleware.IdentityHandler(claims)
		c.Set("JWT_PAYLOAD", claims)
		c.Set("userID", id)

		AuthMiddleware.Authorizator(id, c)

		c.Next()

		return
	}
}

func authRouter(r *gin.Engine) {

	r.POST("/login", AuthMiddleware.LoginHandler)

	authGroup := r.Group("/auth")
	authGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		authGroup.GET("/refresh_token", AuthMiddleware.RefreshHandler)
		authGroup.POST("/telegram", auth.TelegramHandler)
	}
}
