package user

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	telegramUtils "github.com/tokenme/adx/tools/telegram"
	"github.com/tokenme/adx/utils"
	"github.com/tokenme/adx/utils/twilio"
	"github.com/ziutek/mymysql/mysql"
	"net/http"
	"strings"
	"github.com/tokenme/adx/utils/verify253"
)

type CreateRequest struct {
	Mobile       string `form:"mobile" json:"mobile"`
	CountryCode  uint   `form:"country_code" json:"country_code"`
	VerifyCode   string `form:"verify_code" json:"verify_code"`
	Email        string `form:"email" json:"email"`
	Password     string `form:"passwd" json:"passwd" binding:"required"`
	RePassword   string `form:"repasswd" json:"repasswd" binding:"required"`
	Telegram     string `form:"telegram" json:"telegram"`
	IsPublisher  uint   `form:"is_publisher" json:"is_publisher"`
	IsAdvertiser uint   `form:"is_advertiser" json:"is_advertiser"`
}

func CreateHandler(c *gin.Context) {
	var req CreateRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	//log.Info("req: %s", Json(req))

	if Check(req.Email == "" && req.Mobile == "", "missing email and mobile", c) {
		return
	}
	if Check(req.IsPublisher != 1 && req.IsAdvertiser != 1, "missing account type", c) {
		return
	}
	passwdLength := len(req.Password)
	if Check(passwdLength < 6 || passwdLength > 32, "password length must between 6-32", c) {
		return
	}
	if Check(req.Password != req.RePassword, "repassword!=password", c) {
		return
	}
	token, err := uuid.NewV4()
	if CheckErr(err, c) {
		return
	}
	salt := utils.Sha1(token.String())
	token, err = uuid.NewV4()
	if CheckErr(err, c) {
		return
	}
	activationCode := utils.Sha1(token.String())
	passwd := utils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Password, salt))

	mobile := strings.Replace(req.Mobile, " ", "", 0)
	active := 0
	if mobile != "" {
		ret := twilio.AuthVerificationResponse{}
		if req.CountryCode != 86 {
			ret, err = twilio.AuthVerification(Config.TwilioToken, mobile, req.CountryCode, req.VerifyCode)
		}else{
			ret,err = verify253.AuthVerification(mobile,req.VerifyCode,req.CountryCode)
		}
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		if Check(!ret.Success, ret.Message, c) {
			return
		}
		active = 1
	}

	privateKey, _, err := eth.GenerateAccount()
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	walletSalt, wallet, err := utils.AddressEncrypt(privateKey, Config.TokenSalt)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	var telegram common.TelegramUser
	if req.Telegram != "" && telegramUtils.TelegramAuthCheck(req.Telegram, Config.TelegramBotToken) {
		telegram, _ = telegramUtils.ParseTelegramAuth(req.Telegram)
	}
	db := Service.Db
	var res mysql.Result
	if req.Email != "" && req.Mobile != "" {
		_, res, err = db.Query(`INSERT INTO adx.users (country_code, mobile, email, passwd, salt, activation_code, active, telegram_id, telegram_username, telegram_firstname, telegram_lastname, telegram_avatar, is_publisher, is_advertiser) VALUES (%d, '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', '%s', %d, %d)`, req.CountryCode, db.Escape(mobile), db.Escape(req.Email), db.Escape(passwd), db.Escape(salt), db.Escape(activationCode), active, telegram.Id, db.Escape(telegram.Username), db.Escape(telegram.Firstname), db.Escape(telegram.Lastname), db.Escape(telegram.Avatar), req.IsPublisher, req.IsAdvertiser)
	} else if req.Email != "" {
		_, res, err = db.Query(`INSERT INTO adx.users (email, passwd, salt, activation_code, active, telegram_id, telegram_username, telegram_firstname, telegram_lastname, telegram_avatar, is_publisher, is_advertiser) VALUES ('%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', '%s', %d, %d)`, db.Escape(req.Email), db.Escape(passwd), db.Escape(salt), db.Escape(activationCode), active, telegram.Id, db.Escape(telegram.Username), db.Escape(telegram.Firstname), db.Escape(telegram.Lastname), db.Escape(telegram.Avatar), req.IsPublisher, req.IsAdvertiser)
	} else if req.Mobile != "" {
		_, res, err = db.Query(`INSERT INTO adx.users (country_code, mobile, passwd, salt, activation_code, active, telegram_id, telegram_username, telegram_firstname, telegram_lastname, telegram_avatar, is_publisher, is_advertiser) VALUES (%d, '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', '%s', %d, %d)`, req.CountryCode, db.Escape(mobile), db.Escape(passwd), db.Escape(salt), db.Escape(activationCode), active, telegram.Id, db.Escape(telegram.Username), db.Escape(telegram.Firstname), db.Escape(telegram.Lastname), db.Escape(telegram.Avatar), req.IsPublisher, req.IsAdvertiser)
	}
	if err != nil && err.(*mysql.Error).Code == mysql.ER_DUP_ENTRY {
		if req.Email != "" {
			if rows, _, err := db.Query(`SELECT active FROM adx.users WHERE email='%s' LIMIT 1`, db.Escape(req.Email)); err == nil && len(rows) > 0 {
				active := rows[0].Uint(0)
				if active == 0 {
					c.JSON(http.StatusOK, APIError{Code: UNACTIVATED_USER_ERROR, Msg: "account already exists, but did not activate!"})
					return
				}
			}
		}
		c.JSON(http.StatusOK, APIError{Code: DUPLICATE_USER_ERROR, Msg: "account already exists"})
		return
	}
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	userId := res.InsertId()
	_, _, err = db.Query(`INSERT IGNORE INTO adx.user_wallets (user_id, token_type, salt, wallet, name, is_private, is_main) VALUES (%d, 'ETH', '%s', '%s', 'SYS', 1, 1)`, userId, db.Escape(walletSalt), db.Escape(wallet))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if req.Email != "" {
		user := common.User{
			Email:          req.Email,
			ActivationCode: activationCode,
			IsPublisher:    req.IsPublisher,
			IsAdvertiser:   req.IsAdvertiser,
		}
		if err = EmailQueue.NewRegister(user); CheckErr(err, c) {
			raven.CaptureError(err, nil)
		}
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
