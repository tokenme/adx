package media

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type InfoRequest struct {
	Id uint64 `form:"id" json:"id" binding:"required"`
}

func InfoHandler(c *gin.Context) {
	var req InfoRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	db := Service.Db
	var query string
	if user.IsPublisher == 1 && user.IsAdmin == 0{
		query = fmt.Sprintf(`SELECT id, title, domain, intro, salt, verified, online_status, inserted_at, updated_at FROM adx.medias WHERE id=%d AND user_id=%d LIMIT 1`, req.Id, user.Id)
	} else {
		query = fmt.Sprintf(`SELECT id, title, domain, intro, salt, verified, online_status, inserted_at, updated_at FROM adx.medias WHERE id=%d LIMIT 1`, req.Id)
	}
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	if Check(len(rows) == 0, "not found", c) {
		return
	}
	row := rows[0]
	media := common.Media{
		Id:           row.Uint64(0),
		Title:        row.Str(1),
		Domain:       row.Str(2),
		Desc:         row.Str(3),
		Identity:     row.Str(4),
		Verified:     row.Uint(5),
		OnlineStatus: row.Uint(6),
		InsertedAt:   row.ForceLocaltime(7),
		UpdatedAt:    row.ForceLocaltime(8),
	}
	if user.IsPublisher == 1 {
		media = media.Complete()
	}
	c.JSON(http.StatusOK, media)
	return

}
