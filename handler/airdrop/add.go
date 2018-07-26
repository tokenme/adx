package airdrop

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"time"
)

type AddRequest struct {
	Title         string `form:"title" json:"title" binding:"required"`
	TokenAddress  string `form:"token_address" json:"token_address" binding:"required"`
	Budget        uint64 `form:"budget" json:"budget" binding:"required"`
	GiveOut       uint64 `form:"give_out" json:"give_out" binding:"required"`
	Bonus         uint   `form:"bonus" json:"bonus" binding:"required"`
	TelegramGroup string `form:"telegram_group" json:"telegram_group" binding:"required"`
	StartDate     string `form:"start_date" json:"start_date" binding:"required"`
	EndDate       string `form:"end_date" json:"end_date" binding:"required"`
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

	if Check(user.IsAdvertiser != 1, "unauthorized", c) {
		return
	}
	title := utils.Normalize(req.Title)
	telegramGroup := utils.Normalize(req.TelegramGroup)
	db := Service.Db

	rows, _, err := db.Query(`SELECT 1 FROM adx.tokens WHERE address='%s'`, req.TokenAddress)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found the token", c) {
		return
	}
	_, ret, err := db.Query(`INSERT INTO adx.airdrops (user_id, title, token_address, budget, give_out, bonus, telegram_group, start_date, end_date) VALUES (%d, '%s', '%s', %d, %d, '%s', '%s', '%s')`, user.Id, db.Escape(title), db.Escape(req.TokenAddress), req.Budget, req.GiveOut, req.Bonus, db.Escape(telegramGroup), db.Escape(req.StartDate), db.Escape(req.EndDate))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	airdropId := ret.InsertId()
	airdrop := common.Airdrop{
		Id:    airdropId,
		Title: title,
		Token: common.Token{
			Address: req.TokenAddress,
		},
		Budget:        req.Budget,
		GiveOut:       req.GiveOut,
		Bonus:         req.Bonus,
		TelegramGroup: telegramGroup,
		InsertedAt:    time.Now(),
		UpdatedAt:     time.Now(),
	}
	c.JSON(http.StatusOK, airdrop)
}
