package auth

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils/twilio"
	"net/http"
	"strings"
)

type VerifyRequest struct {
	Mobile  string `form:"mobile" json:"mobile" binding:"required"`
	Code    string `form:"code" json:"code" binding:"required"`
	Country uint   `form:"country" json:"country" binding:"required"`
}

func VerifyHandler(c *gin.Context) {
	var req VerifyRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	mobile := strings.Replace(req.Mobile, " ", "", 0)
	db := Service.Db
	rows, _, err := db.Query(`SELECT 1 FROM adx.auth_verify_codes WHERE country_code=%d AND mobile='%s' AND code='%s' LIMIT 1`, req.Country, db.Escape(mobile), db.Escape(req.Code))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if len(rows) > 0 {
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}
	ret, err := twilio.AuthVerification(Config.TwilioToken, req.Mobile, req.Country, req.Code)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(!ret.Success, ret.Message, c) {
		return
	}
	_, _, err = db.Query(`DELETE FROM adx.auth_verify_codes WHERE country_code=%d AND mobile='%s'`, req.Country, db.Escape(mobile))
	_, _, err = db.Query(`INSERT INTO adx.auth_verify_codes (country_code, mobile, code) VALUES (%d, '%s', '%s') ON DUPLICATE KEY UPDATE inserted=NOW()`, req.Country, db.Escape(mobile), db.Escape(req.Code))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
