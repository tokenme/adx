package media

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
	Id              uint64 `form:"id" json:"id"`
	Title           string `form:"title" json:"title"`
	placeholder_img string `from:"placeholder_img" json:"placeholder_img"`
	OnlineStatus    uint   `form:"online_status" json:"online_status"`
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

	if Check(user.IsPublisher != 1 && user.IsAdmin != 1, "unauthorized", c) {
		return
	}

	db := Service.Db

	title := utils.Normalize(req.Title)
	url := utils.Normalize(req.placeholder_img)
	var set = []string{fmt.Sprintf("online_status=%d", req.OnlineStatus)}
	if title != "" {
		set = append(set, fmt.Sprintf("title='%s'", db.Escape(title)))
	}
	if url != "" {
		set = append(set, fmt.Sprintf("url='%s'", db.Escape(url)))
	}
	var err error
	if user.IsAdmin == 1 {
		_, _, err = db.Query(`UPDATE adx.medias SET %s WHERE id=%d `, strings.Join(set, ","), req.Id)

	} else {
		_, _, err = db.Query(`UPDATE adx.medias SET %s WHERE id=%d AND user_id=%d`, strings.Join(set, ","), req.Id, user.Id)
	}
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
