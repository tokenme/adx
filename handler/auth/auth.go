package auth

import (
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/coins/eth"
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
	} else if loginInfo.Mobile != "" && loginInfo.Password != "" {
		if loginInfo.IsPublisher == 1 {
			where = fmt.Sprintf("country_code=%d AND mobile='%s' AND is_publisher=1", loginInfo.CountryCode, db.Escape(loginInfo.Mobile))
		} else if loginInfo.IsAdvertiser == 1 {
			where = fmt.Sprintf("country_code=%d AND mobile='%s' AND is_advertiser=1", loginInfo.CountryCode, db.Escape(loginInfo.Mobile))
		} else if loginInfo.IsAdmin == 1 {
			where = fmt.Sprintf("country_code=%d AND mobile='%s' AND is_admin=1", loginInfo.CountryCode, db.Escape(loginInfo.Mobile))
		}
	}
	if where == "" {
		return loginInfo.Email, false
	}
	query := `SELECT 
                u.id, 
			     u.country_code,
			     u.mobile,
                u.email, 
                u.salt, 
                u.passwd,
                u.telegram_id,
                u.telegram_username,
                u.telegram_firstname,
                u.telegram_lastname,
                u.telegram_avatar,
                uw.salt AS uw_salt,
                uw.wallet
            FROM adx.users AS u
            INNER JOIN adx.user_wallets AS uw ON (uw.user_id = u.id AND uw.is_main = 1 AND uw.token_type='ETH')
            WHERE %s
            AND u.active = 1
            LIMIT 1`
	rows, res, err := db.Query(query, where)
	if err != nil || len(rows) == 0 {
		return loginInfo.Email, false
	}
	row := rows[0]
	user := common.User{
		Id:          row.Uint64(res.Map("id")),
		CountryCode: row.Uint(res.Map("country_code")),
		Mobile:      row.Str(res.Map("mobile")),
		Email:       row.Str(res.Map("email")),
		Salt:        row.Str(res.Map("salt")),
		Password:    row.Str(res.Map("passwd")),
	}
	telegramId := row.Int64(res.Map("telegram_id"))
	if telegramId > 0 {
		telegram := &common.TelegramUser{
			Id:        telegramId,
			Username:  row.Str(res.Map("telegram_username")),
			Firstname: row.Str(res.Map("telegram_firstname")),
			Lastname:  row.Str(res.Map("telegram_lastname")),
			Avatar:    row.Str(res.Map("telegram_avatar")),
		}
		user.Telegram = telegram
	}
	walletSalt := row.Str(res.Map("uw_salt"))
	walletEncrypt := row.Str(res.Map("wallet"))
	privateKey, err := utils.AddressDecrypt(walletEncrypt, walletSalt, Config.TokenSalt)
	if err != nil {
		return loginInfo.Email, false
	}
	publicKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if err != nil {
		return loginInfo.Email, false
	}
	user.Wallet = publicKey
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
