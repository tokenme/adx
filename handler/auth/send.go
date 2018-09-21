package auth

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils/twilio"
	"github.com/tokenme/adx/utils/verify253"
	"net/http"
	"strings"
	"fmt"
)

type SendRequest struct {
	Mobile  string `form:"mobile" json:"mobile" binding:"required"`
	Country uint   `form:"country" json:"country" binding:"required"`
}

func SendHandler(c *gin.Context) {
	var req SendRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	fmt.Printf("req: %s\n", Json(req))
	mobile := strings.Replace(req.Mobile, " ", "", 0)
	if req.Country == 86 {
		verify253.AuthSend(mobile, c)
		return
	}
	ret, err := twilio.AuthSend(Config.TwilioToken, mobile, req.Country)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(!ret.Success, ret.Message, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
