package user

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
)

type ResetPasswordRequest struct {
	Code       string `form:"code" json:"code" binding:"required"`
	Password   string `form:"passwd" json:"passwd" binding:"required"`
	RePassword string `form:"repasswd" json:"repasswd" binding:"required"`
}

func ResetPasswordHandler(c *gin.Context) {
	var req ResetPasswordRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	passwdLength := len(req.Password)
	if Check(passwdLength < 6 || passwdLength > 32, "password length must between 6-32", c) {
		return
	}
	if Check(req.Password != req.RePassword, "repassword!=password", c) {
		return
	}

	db := Service.Db
	rows, _, err := db.Query(`SELECT id, salt FROM adx.users WHERE reset_pwd_code='%s' LIMIT 1`, db.Escape(req.Code))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(len(rows) == 0, "user doesn't exists", c) {
		return
	}
	userId := rows[0].Uint64(0)
	salt := rows[0].Str(1)
	passwd := utils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Password, salt))
	_, _, err = db.Query(`UPDATE adx.users SET passwd='%s', reset_pwd_code=NULL WHERE id=%d`, db.Escape(passwd), userId)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
