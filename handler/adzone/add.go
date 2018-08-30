package adzone

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"github.com/ziutek/mymysql/mysql"
	"github.com/pkg/errors"
)

type AddRequest struct {
	MediaId        uint64  `form:"media_id" json:"media_id" binding:"required"`
	SizeId         uint    `form:"size_id" json:"size_id" binding:"required"`
	MinCPT         float64 `form:"min_cpt" json:"min_cpt" binding:"required"`
	Rolling        uint    `form:"rolling" json:"rolling" binding:"required"`
	Settlement     uint    `form:"settlement" json:"settlement" binding:"required"`
	Url            string  `form:"url" json:"url" binding:"required"`
	Desc           string  `form:"desc" json:"desc" binding:"required"`
	PlaceholderUrl string  `form:"placeholder_url" json:"placeholder_url"`
	PlaceholderImg string  `form:"placeholder_img" json:"placeholder_img"`
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

	if Check(user.IsPublisher != 1&& user.IsAdmin !=1, "unauthorized", c) {
		return
	}
	db := Service.Db
	rows:=[]mysql.Row{}
	err:=errors.New("")
	if user.IsAdmin == 1{
		rows, _, err = db.Query(`SELECT user_id FROM adx.medias WHERE id=%d LIMIT 1`, req.MediaId)

	}else {
		rows, _, err = db.Query(`SELECT user_id FROM adx.medias WHERE id=%d AND user_id=%d LIMIT 1`, req.MediaId, user.Id)
	}
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(len(rows) == 0, "media not found", c) {
		return
	}
	mediaUserId := rows[0].Uint64(0)
	desc := utils.Normalize(req.Desc)
	if req.PlaceholderImg != "" && req.PlaceholderUrl != "" {
		_, _, err = db.Query(`INSERT INTO adx.adzones (user_id, media_id, size_id, min_cpt, settlement, rolling, url, intro, placeholder_url, placeholder_img) VALUES (%d, %d, %d, %.18f, %d, %d, '%s', '%s', '%s', '%s')`, mediaUserId, req.MediaId, req.SizeId, req.MinCPT, req.Settlement, req.Rolling, db.Escape(req.Url), db.Escape(desc), db.Escape(req.PlaceholderUrl), db.Escape(req.PlaceholderImg))
	} else {
		_, _, err = db.Query(`INSERT INTO adx.adzones (user_id, media_id, size_id, min_cpt, settlement, rolling, url, intro) VALUES (%d, %d, %d, %.18f, %d, %d, '%s', '%s')`, mediaUserId, req.MediaId, req.SizeId, req.MinCPT, req.Settlement, req.Rolling, db.Escape(req.Url), db.Escape(desc))
	}

	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
