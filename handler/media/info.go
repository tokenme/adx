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
	if user.IsPublisher == 1 && user.IsAdmin == 0 {
		query = fmt.Sprintf(`SELECT a.id, a.title, a.domain, a.url, a.salt, a.verified, a.online_status, a.inserted_at, 
		a.updated_at FROM adx.medias AS a 
		WHERE a.id=%d AND a.user_id=%d LIMIT 1`, req.Id, user.Id)
	} else {
		query = fmt.Sprintf(`SELECT a.id, a.title, a.domain, 
		a.url, a.salt, a.verified, a.online_status, a.inserted_at, 
		a.updated_at 
		FROM adx.medias AS a
		WHERE a.id=%d LIMIT 1`, req.Id)
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
	ImgUrl := ""
	Img := &common.PrivateAuctionCreative{
		Url: row.Str(2),
		Img: row.Str(3)}
	if Img.Url != "" && Img.Img != "" {
		ImgUrl = Img.GetImgUrl(Config)
	}
	media := common.Media{
		Id:           row.Uint64(0),
		Title:        row.Str(1),
		Domain:       row.Str(2),
		ImgUrl:       ImgUrl,
		Identity:     row.Str(4),
		Verified:     row.Uint(5),
		OnlineStatus: row.Uint(6),
		InsertedAt:   row.ForceLocaltime(7),
		UpdatedAt:    row.ForceLocaltime(8),
	}
	query = `SELECT adzone.id,
adzone.size_id,size.width,
size.height,media.id AS media_id,adzone.intro,
adzone.min_cpt,adzone.min_cpm,
adzone.url,
adzone.user_id,
adzone.rolling,adzone.inserted_at,
adzone.updated_at,
adzone.settlement,
adzone.online_status,
adzone.advantage,
adzone.location,
adzone.traffic
FROM adx.adzones AS adzone
INNER JOIN adx.medias AS media ON (media.id = adzone.media_id)
INNER JOIN adx.sizes AS size   ON (adzone.size_id = size.id)
WHERE media.id = %d AND adzone.online_status = 1 GROUP BY adzone.id`
	rows, Result, err := db.Query(query, req.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	Adzones := []common.Adzone{}
	for _, row := range rows {
		Adzones = append(Adzones, common.Adzone{
			Id:     row.Uint64(Result.Map(`id`)),
			UserId: row.Uint64(Result.Map(`user_id`)),
			Media:  media,
			Size: common.Size{
				Id:     row.Uint(Result.Map(`size_id`)),
				Width:  row.Uint(Result.Map(`width`)),
				Height: row.Uint(Result.Map(`height`)),
			},
			Url:          row.Str(Result.Map(`url`)),
			MinCPT:       row.Float(Result.Map(`min_cpt`)),
			MinCPM:       row.Float(Result.Map(`min_cpm`)),
			Intro:        row.Str(Result.Map(`intro`)),
			Rolling:      row.Uint(Result.Map(`rolling`)),
			OnlineStatus: row.Uint(Result.Map(`online_status`)),
			InsertedAt:   row.ForceLocaltime(Result.Map(`inserted_at`)),
			UpdatedAt:    row.ForceLocaltime(Result.Map(`updated_at`)),
			Advantage:    row.Str(Result.Map("advantage")),
			Location:     row.Str(Result.Map("location")),
			Traffic:      row.Str(Result.Map("traffic")),
		})
	}
	media.Adzones = Adzones
	if user.IsPublisher == 1 {
		media = media.Complete()
	}
	c.JSON(http.StatusOK, media)
	return
}
