package airdropadzone

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func ListGetHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	channelIdStr := c.Query("channel_id")
	channelId, err := Uint64NonZero(channelIdStr, "missing channel_id")
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT id, user_id, channel_id, name FROM adx.airdrop_adzones WHERE user_id=%d AND channel_id=%d ORDER BY id ASC`, user.Id, channelId)
	if CheckErr(err, c) {
		return
	}
	var adzones []common.AirdropAdzone
	for _, row := range rows {
		adzones = append(adzones, common.AirdropAdzone{
			Id:        row.Uint64(0),
			UserId:    row.Uint64(1),
			ChannelId: row.Uint64(2),
			Name:      row.Str(3),
		})
	}
	c.JSON(http.StatusOK, adzones)
}
