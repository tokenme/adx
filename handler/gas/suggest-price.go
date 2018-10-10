package gas

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"math/big"
	"net/http"
)

func SuggestPriceHandler(c *gin.Context) {
	geth := Service.Geth
	if geth == nil {
		c.JSON(http.StatusOK, APIError{Code: NOTFOUND_ERROR, Msg: "not found"})
		return
	}
	price, err := geth.SuggestGasPrice(c)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, map[string]*big.Int{"price": price})
}
