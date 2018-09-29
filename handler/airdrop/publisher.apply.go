package airdrop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nlopes/slack"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/ziutek/mymysql/mysql"
	"net/http"
)

func PublisherApplyHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	_, _, err := db.Query(`INSERT IGNORE INTO adx.airdrop_applications (user_id) VALUES (%d)`, user.Id)
	if Check(err != nil && err.(*mysql.Error).Code == mysql.ER_DUP_ENTRY, "Already applied", c) {
		return
	}
	if CheckErr(err, c) {
		return
	}
	if Service.Slack != nil {
		params := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color:      "#1976d2",
			AuthorName: user.ShowName,
			AuthorIcon: user.Avatar,
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Mobile",
					Value: user.Mobile,
					Short: true,
				},
				slack.AttachmentField{
					Title: "CountryCode",
					Value: fmt.Sprintf("%d", user.CountryCode),
					Short: true,
				},
			},
		}
		params.Attachments = []slack.Attachment{attachment}
		_, _, err = Service.Slack.PostMessage("G9Y7METUG", "new user applied for publisher", params)
		if err != nil {
			log.Error(err.Error())
		}
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
