package airdrop

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"strings"
	"time"
)

type UpdateRequest struct {
	Id             uint64 `form:"id" json:"id" binding:"required"`
	GasPrice       uint64 `form:"gas_price" json:"gas_price"`
	GasLimit       uint64 `form:"gas_limit" json:"gas_limit"`
	GiveOut        uint64 `form:"give_out" json:"give_out"`
	DropDate       int64  `form:"drop_date" json:"drop_date"`
	MaxSubmissions int    `form:"max_submissions" json:"max_submissions"`
	ReplyMsg       string `form:"reply_msg" json:"reply_msg"`
	Status         uint   `form:"status" json:"status"`
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
	if Check(user.IsAirdropPublisher == 0 && user.IsAdmin == 0, "invalid permission", c) {
		return
	}
	db := Service.Db
	var updateFields []string
	if req.GasPrice > 0 {
		updateFields = append(updateFields, fmt.Sprintf("gas_price=%d", req.GasPrice))
	}
	if req.GasLimit > 0 {
		updateFields = append(updateFields, fmt.Sprintf("gas_limit=%d", req.GasLimit))
	}
	if req.GiveOut > 0 {
		updateFields = append(updateFields, fmt.Sprintf("give_out=%d", req.GiveOut))
	}
	if req.DropDate > 0 {
		updateFields = append(updateFields, fmt.Sprintf("drop_date='%s'", time.Unix(req.DropDate/1000, 0).Format("2006-01-02")))
	}
	replyMsg := "NULL"
	if req.ReplyMsg != "" {
		replyMsg = fmt.Sprintf("'%s'", db.Escape(req.ReplyMsg))
	}
	updateFields = append(updateFields, fmt.Sprintf("reply_msg=%s", replyMsg))
	updateFields = append(updateFields, fmt.Sprintf("max_submissions=%d", req.MaxSubmissions))
	updateFields = append(updateFields, fmt.Sprintf("status=%d", req.Status))
	if len(updateFields) == 0 {
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}

	var checkUser string
	if user.IsAdmin == 0 {
		checkUser = fmt.Sprintf(" AND user_id=%d", user.Id)
	}
	_, _, err := db.Query(`UPDATE adx.airdrops SET %s WHERE id=%d%s LIMIT 1`, strings.Join(updateFields, ","), req.Id, checkUser)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
