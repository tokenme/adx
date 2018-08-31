package adzone

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
	Id             uint64  `form:"id" json:"id"`
	Url            string  `form:"url" json:"url"`
	Desc           string  `form:"desc" json:"desc"`
	Rolling        uint    `form:"rolling" json:"rolling"`
	MinCPT         float64 `from:"min_cpt" json:"min_cpt"`
	OnlineStatus   uint    `form:"online_status" json:"online_status"`
	PlaceholderUrl string  `form:"placeholder_url" json:"placeholder_url"`
	PlaceholderImg string  `form:"placeholder_img" json:"placeholder_img"`
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

	desc := utils.Normalize(req.Desc)
	var set = []string{fmt.Sprintf("online_status=%d", req.OnlineStatus)}
	if req.Rolling > 0 {
		set = append(set, fmt.Sprintf("rolling=%d", req.Rolling))
	}
	if req.MinCPT > 0 {
		set = append(set, fmt.Sprintf("min_cpt=%.18f", req.MinCPT))
	}
	if desc != "" {
		set = append(set, fmt.Sprintf("intro='%s'", db.Escape(desc)))
	}

	if req.Url != "" {
		set = append(set, fmt.Sprintf("url='%s'", db.Escape(req.Url)))
	}

	if req.PlaceholderUrl != "" && req.PlaceholderImg != "" {
		set = append(set, fmt.Sprintf("placeholder_url='%s', placeholder_img='%s'", db.Escape(req.PlaceholderUrl), db.Escape(req.PlaceholderImg)))
	}
	var err error
	if user.IsAdmin == 1 {
		_, _, err = db.Query(`UPDATE adx.adzones SET %s WHERE id=%d `, strings.Join(set, ","), req.Id)
	} else {
		_, _, err = db.Query(`UPDATE adx.adzones SET %s WHERE id=%d AND user_id=%d`, strings.Join(set, ","), req.Id, user.Id)

	}
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
