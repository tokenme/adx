package user

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func BalanceHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	balance, err := user.Balance(c, Service, Config)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, balance)
}
