package media

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func ListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if Check(user.IsPublisher != 1 && user.IsAdvertiser !=1 &&user.IsAdmin !=1, "unauthorized", c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT id, title, domain, intro, salt, verified, online_status, inserted_at, updated_at FROM adx.medias WHERE user_id=%d ORDER BY id DESC`, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	var medias []common.Media
	for _, row := range rows {
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
		media = media.Complete()
		medias = append(medias, media)
	}

	c.JSON(http.StatusOK, medias)
	return

}
