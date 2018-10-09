package airdropadzone

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type AddRequest struct {
	ChannelId uint64 `form:"channel_id" json:"channel_id" binding:"required"`
	Name      string `form:"name" json:"name" binding:"required"`
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
	if !userOwnsChannel(user.Id, req.ChannelId) {
		CheckErr(fmt.Errorf("channel doesn't belong to user"), c)
	}
	db := Service.Db
	_, _, err := db.Query(`INSERT IGNORE INTO adx.airdrop_adzones (channel_id, user_id, name) VALUES (%d, %d, '%s')`, req.ChannelId, user.Id, db.Escape(req.Name))
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}

func userOwnsChannel(userId uint64, channelId uint64) bool {
	db := Service.Db
	rows, _, err := db.Query(`SELECT 1 FROM adx.channels WHERE user_id IN (0, %d) AND id=%d LIMIT 1`, userId, channelId)
	if len(rows) == 0 || err != nil {
		return false
	}
	return true
}
