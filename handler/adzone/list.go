package adzone

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type ListRequest struct {
	MediaId uint64 `form:"media_id" json:"media_id" binding:"required"`
}

func ListHandler(c *gin.Context) {
	var req ListRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if Check(user.IsPublisher != 1, "unauthorized", c) {
		return
	}

	db := Service.Db
	rows, _, err := db.Query(`SELECT a.id, a.url, a.size_id, s.width, s.height, a.min_cpm, a.min_cpt, a.settlement, a.rolling, a.intro, a.online_status, m.id, m.title, m.domain, a.inserted_at, a.updated_at FROM adx.adzones AS a INNER JOIN adx.medias AS m ON (m.id=a.media_id) INNER JOIN adx.sizes AS s ON (s.id=a.size_id) WHERE a.media_id=%d AND a.user_id=%d ORDER BY a.id DESC`, req.MediaId, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	var adzones []common.Adzone
	for _, row := range rows {
		adzone := common.Adzone{
			Id:  row.Uint64(0),
			Url: row.Str(1),
			Size: common.Size{
				Id:     row.Uint(2),
				Width:  row.Uint(3),
				Height: row.Uint(4),
			},
			MinCPM:       row.ForceFloat(5),
			MinCPT:       row.ForceFloat(6),
			Settlement:   row.Uint(7),
			Rolling:      row.Uint(8),
			Desc:         row.Str(9),
			OnlineStatus: row.Uint(10),
			Media: common.Media{
				Id:     row.Uint64(11),
				Title:  row.Str(12),
				Domain: row.Str(13),
			},
			InsertedAt: row.ForceLocaltime(14),
			UpdatedAt:  row.ForceLocaltime(15),
		}
		adzones = append(adzones, adzone)
	}

	c.JSON(http.StatusOK, adzones)
	return

}
