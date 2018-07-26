package airdrop

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"strings"
)

type UpdateRequest struct {
	Id            uint64 `form:"id" json:"id" binding:"required"`
	Title         string `form:"title" json:"title"`
	Budget        uint64 `form:"budget" json:"budget"`
	GiveOut       uint64 `form:"give_out" json:"give_out"`
	Bonus         uint   `form:"bonus" json:"bonus"`
	TelegramGroup string `form:"telegram_group" json:"telegram_group"`
	StartDate     string `form:"start_date" json:"start_date"`
	EndDate       string `form:"end_date" json:"end_date"`
	OnlineStatus  int    `form:"online_status" json:"online_status"`
}

func UpdateHandler(c *gin.Context) {
	var req UpdateRequest
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

	rows, _, err := db.Query(`SELECT 1 FROM adx.airdrops WHERE id=%d AND user_id=%d AND start_date>DATE(NOW()) LIMIT 1`, req.Id, user.Id)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not allowed to edit", c) {
		return
	}
	var sets []string
	if title != "" {
		sets = append(sets, fmt.Sprintf("title='%s'", db.Escape(title)))
	}
	if req.Budget > 0 {
		sets = append(sets, fmt.Sprintf("budget=%d", req.Budget))
	}
	if req.GiveOut > 0 {
		sets = append(sets, fmt.Sprintf("give_out=%d", req.GiveOut))
	}
	if req.Bonus > 0 {
		sets = append(sets, fmt.Sprintf("bonus=%d", req.Bonus))
	}
	if telegramGroup != "" {
		sets = append(sets, fmt.Sprintf("telegram_group='%s'", db.Escape(telegramGroup)))
	}
	if req.StartDate != "" {
		sets = append(sets, fmt.Sprintf("start_date='%s'", db.Escape(req.StartDate)))
	}
	if req.EndDate != "" {
		sets = append(sets, fmt.Sprintf("end_date='%s'", db.Escape(req.EndDate)))
	}
	if req.OnlineStatus == 1 || req.OnlineStatus == -1 {
		sets = append(sets, fmt.Sprintf("online_status=%d", req.OnlineStatus))
	}
	if Check(len(sets) == 0, "nothing to update", c) {
		return
	}
	_, _, err = db.Query(`UPDATE adx.airdrops SET %s WHERE id=%d AND user_id=%d`, strings.Join(sets, ","), req.Id, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
