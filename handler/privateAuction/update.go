package privateAuction

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
	Id           uint64  `form:"id" json:"id" binding:"required"`
	Title        string  `form:"title" json:"title"`
	Price        float64 `from:"price" json:"price"`
	OnlineStatus uint    `form:"online_status" json:"online_status"`
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

	db := Service.Db

	title := utils.Normalize(req.Title)
	var set = []string{fmt.Sprintf("online_status=%d", req.OnlineStatus)}
	if req.Price > 0 {
		set = append(set, fmt.Sprintf("price=%.18f", req.Price))
	}
	if title != "" {
		set = append(set, fmt.Sprintf("title='%s'", db.Escape(title)))
	}

	_, _, err := db.Query(`UPDATE adx.private_auctions SET %s WHERE id=%d AND user_id=%d AND audit_status=0`, strings.Join(set, ","), req.Id, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
