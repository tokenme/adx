package auth

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	telegramUtils "github.com/tokenme/adx/tools/telegram"
	"net/http"
)

type TelegramRequest struct {
	Telegram string `form:"telegram" json:"telegram" binding:"required"`
}

func TelegramHandler(c *gin.Context) {
	var req TelegramRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(!telegramUtils.TelegramAuthCheck(req.Telegram, Config.TelegramBotToken), "invalid hash", c) {
		return
	}
	telegram, err := telegramUtils.ParseTelegramAuth(req.Telegram)
	if CheckErr(err, c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	db := Service.Db
	_, _, err = db.Query(`UPDATE adx.users SET telegram_id=%d, telegram_username='%s', telegram_firstname='%s', telegram_lastname='%s', telegram_avatar='%s' WHERE id=%d`, telegram.Id, db.Escape(telegram.Username), db.Escape(telegram.Firstname), db.Escape(telegram.Lastname), db.Escape(telegram.Avatar), user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
