package airdrop

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"time"
)

type AddRequest struct {
	Title          string `form:"title" json:"title" binding:"required"`
	TokenAddress   string `form:"token_address" json:"token_address" binding:"required"`
	GasPrice       uint64 `form:"gas_price" json:"gas_price" binding:"required"`
	GasLimit       uint64 `form:"gas_limit" json:"gas_limit"`
	GiveOut        uint64 `form:"give_out" json:"give_out" binding:"required"`
	Bonus          uint   `form:"bonus" json:"bonus" binding:"required"`
	RequireEmail   uint   `form:"require_email" json:"require_email" binding:"required"`
	TelegramGroup  string `form:"telegram_group" json:"telegram_group" binding:"required"`
	ReplyMsg       string `form:"reply_msg" json:"reply_msg"`
	StartDate      int64  `form:"start_date" json:"start_date" binding:"required"`
	EndDate        int64  `form:"end_date" json:"end_date" binding:"required"`
	MaxSubmissions uint   `form:"max_submissions json:"max_submissions"`
}

func AddHandler(c *gin.Context) {
	var req AddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsPublisher == 0 && user.IsAdmin == 0, "invalid permission", c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT address, name, symbol, decimals, protocol FROM adx.tokens WHERE address='%s' LIMIT 1`, db.Escape(req.TokenAddress))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "missing token", c) {
		return
	}
	token := common.Token{
		Address:  rows[0].Str(0),
		Name:     rows[0].Str(1),
		Symbol:   rows[0].Str(2),
		Decimals: rows[0].Uint(3),
		Protocol: rows[0].Str(4),
	}
	privateKey, publicKey, err := eth.GenerateAccount()
	if CheckErr(err, c) {
		return
	}
	salt, wallet, err := utils.AddressEncrypt(privateKey, Config.TokenSalt)
	if CheckErr(err, c) {
		return
	}
	if req.GasLimit == 0 {
		req.GasLimit = 420000
	}
	now := time.Now()
	airdrop := common.Airdrop{
		User:           common.User{Id: user.Id, Email: user.Email},
		Title:          req.Title,
		Wallet:         publicKey,
		WalletPrivKey:  privateKey,
		Token:          token,
		GasPrice:       req.GasPrice,
		GasLimit:       req.GasLimit,
		GiveOut:        req.GiveOut,
		Bonus:          req.Bonus,
		RequireEmail:   req.RequireEmail,
		TelegramGroup:  req.TelegramGroup,
		StartDate:      time.Unix(req.StartDate/1000, 0),
		EndDate:        time.Unix(req.EndDate/1000, 0),
		Status:         common.AirdropStatusStop,
		BalanceStatus:  common.AirdropBalanceStatusEmpty,
		CommissionFee:  Config.Airdrop.CommissionFee,
		MaxSubmissions: req.MaxSubmissions,
		ReplyMsg:       req.ReplyMsg,
		Inserted:       now,
		Updated:        now,
	}
	airdrop.DropDate = airdrop.EndDate.Add(24 * time.Hour)
	replyMsg := "NULL"
	if airdrop.ReplyMsg != "" {
		replyMsg = fmt.Sprintf("'%s'", db.Escape(airdrop.ReplyMsg))
	}
	_, ret, err := db.Query(`INSERT INTO adx.airdrops (user_id, title, wallet, salt, token_address, gas_price, gas_limit, commission_fee, give_out, bonus, start_date, end_date, drop_date, telegram_group, require_email, max_submissions, reply_msg) VALUES (%d, '%s', '%s', '%s', '%s', %d, %d, %d, %d, %d, '%s', '%s', '%s', '%s', %d, %d, %s)`, user.Id, db.Escape(airdrop.Title), db.Escape(wallet), db.Escape(salt), db.Escape(token.Address), airdrop.GasPrice, airdrop.GasLimit, airdrop.CommissionFee, airdrop.GiveOut, airdrop.Bonus, db.Escape(airdrop.StartDate.Format("2006-01-02")), db.Escape(airdrop.EndDate.Format("2006-01-02")), db.Escape(airdrop.DropDate.Format("2006-01-02")), db.Escape(airdrop.TelegramGroup), airdrop.RequireEmail, airdrop.MaxSubmissions, replyMsg)
	if CheckErr(err, c) {
		return
	}
	airdrop.Id = ret.InsertId()
	c.JSON(http.StatusOK, airdrop)
}
