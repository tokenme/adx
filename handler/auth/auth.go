package auth

import (
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/middlewares/jwt"
	telegramUtils "github.com/tokenme/adx/tools/telegram"
	"github.com/tokenme/adx/utils"
)

var AuthenticatorFunc = func(loginInfo jwt.Login, c *gin.Context) (string, bool) {
	db := Service.Db
	var where string
	if loginInfo.Telegram != "" {
		if !telegramUtils.TelegramAuthCheck(loginInfo.Telegram, Config.TelegramBotToken) {
			log.Error("Wrong checksum")
			return loginInfo.Email, false
		}
		telegram, err := telegramUtils.ParseTelegramAuth(loginInfo.Telegram)
		if err != nil {
			raven.CaptureError(err, nil)
			return loginInfo.Email, false
		}
		where = fmt.Sprintf("telegram_id=%d", telegram.Id)
	} else if loginInfo.Email != "" && loginInfo.Password != "" {
		if loginInfo.IsPublisher == 1 {
			where = fmt.Sprintf("email='%s' AND is_publisher=1", db.Escape(loginInfo.Email))
		} else if loginInfo.IsAdvertiser == 1 {
			where = fmt.Sprintf("email='%s' AND is_advertiser=1", db.Escape(loginInfo.Email))
		} else if loginInfo.IsAdmin == 1 {
			where = fmt.Sprintf("email='%s' AND is_admin=1", db.Escape(loginInfo.Email))
		}
	}
	if where == "" {
		return loginInfo.Email, false
	}
	query := `SELECT 
                id, 
                email, 
                salt, 
                passwd,
                telegram_id,
                telegram_username,
                telegram_firstname,
                telegram_lastname,
                telegram_avatar
            FROM adx.users
            WHERE %s
            AND active = 1
            LIMIT 1`
	rows, _, err := db.Query(query, where)
	if err != nil || len(rows) == 0 {
		return loginInfo.Email, false
	}
	row := rows[0]
	user := common.User{
		Id:       row.Uint64(0),
		Email:    row.Str(1),
		Salt:     row.Str(2),
		Password: row.Str(3),
	}
	telegramId := row.Int64(4)
	if telegramId > 0 {
		telegram := &common.TelegramUser{
			Id:        telegramId,
			Username:  row.Str(5),
			Firstname: row.Str(6),
			Lastname:  row.Str(7),
			Avatar:    row.Str(8),
		}
		user.Telegram = telegram
	}
	if loginInfo.IsPublisher == 1 {
		user.IsPublisher = 1
	} else if loginInfo.IsAdvertiser == 1 {
		user.IsAdvertiser = 1
	} else if loginInfo.IsAdmin == 1 {
		user.IsAdmin = 1
	}
	user.ShowName = user.GetShowName()
	user.Avatar = user.GetAvatar(Config.CDNUrl)
	c.Set("USER", user)
	passwdSha1 := utils.Sha1(fmt.Sprintf("%s%s%s", user.Salt, loginInfo.Password, user.Salt))
	userPassword := user.Password
	user.Password = ""
	js, err := json.Marshal(user)
	if err != nil {
		return loginInfo.Email, false
	}
	return string(js), passwdSha1 == userPassword || loginInfo.Telegram != ""
}

var AuthorizatorFunc = func(data string, c *gin.Context) bool {
	var user common.User
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return false
	}
	db := Service.Db
	query := `SELECT 1 FROM adx.users WHERE id=%d AND active=1`
	if user.IsPublisher == 1 {
		query = fmt.Sprintf("%s AND is_publisher=1 LIMIT 1", query)
	} else if user.IsAdvertiser == 1 {
		query = fmt.Sprintf("%s AND is_advertiser=1 LIMIT 1", query)
	} else if user.IsAdmin == 1 {
		query = fmt.Sprintf("%s AND is_admin=1 LIMIT 1", query)
	}
	rows, _, err := db.Query(query, user.Id)
	if err != nil || len(rows) == 0 {
		if err != nil {
			raven.CaptureError(err, nil)
		}
		return false
	}
	c.Set("USER", user)
	return true
}
