package router

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/handler/promotion"
	"net/http"
	"time"
)

func promotionRouter(r *gin.Engine) {

	limiter := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limiter.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	limiterHandler := tollbooth_gin.LimitHandler(limiter)

	promotionGroup := r.Group("/promotion")
	promotionGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		promotionGroup.POST("/add", promotion.AddHandler)
		promotionGroup.GET("/list", promotion.ListHandler)
		promotionGroup.GET("/get", promotion.GetHandler)
		promotionGroup.GET("/stats", promotion.StatsHandler)
	}
	r.GET("/promotion/wallet", limiterHandler, promotion.NewWalletHandler)
	r.GET("/promotion/show/:key", limiterHandler, promotion.ShowHandler)
	r.POST("/promotion/submit", limiterHandler, promotion.SubmitHandler)
	r.GET("/promo/:key", limiterHandler, func(c *gin.Context) {
		key := c.Param("key")
		c.Redirect(http.StatusFound, fmt.Sprintf("/promo.html#/%s", key))
	})

	r.GET("/promotion/topn", promotion.TopNHandler)
}
