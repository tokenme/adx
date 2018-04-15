package user

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
)

type SendResetPasswordEmailRequest struct {
	Email        string `form:"email" json:"email" binding:"required"`
	IsPublisher  uint   `form:"is_publisher" json:"is_publisher"`
	IsAdvertiser uint   `form:"is_advertiser" json:"is_advertiser"`
}

func SendResetPasswordEmailHandler(c *gin.Context) {
	var req SendResetPasswordEmailRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.IsPublisher != 1 && req.IsAdvertiser != 1, "missing account type", c) {
		return
	}
	token, err := uuid.NewV4()
	if CheckErr(err, c) {
		return
	}
	code := utils.Sha1(token.String())
	db := Service.Db
	_, _, err = db.Query(`UPDATE adx.users SET reset_pwd_code='%s', reset_pwd_time=NOW() WHERE email='%s'`, db.Escape(code), db.Escape(req.Email))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	user := common.User{
		Email:        req.Email,
		ResetPwdCode: code,
		IsPublisher:  req.IsPublisher,
		IsAdvertiser: req.IsAdvertiser,
	}
	err = EmailQueue.NewResetPwd(user)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
