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
	rows, _, err := db.Query(`SELECT a.id, a.url, a.size_id, s.width, s.height, a.min_cpm, a.min_cpt, a.settlement, a.rolling, a.intro, a.online_status, a.placeholder_img, a.placeholder_url, m.id, m.title, m.domain, m.online_status, a.inserted_at, a.updated_at FROM adx.adzones AS a INNER JOIN adx.medias AS m ON (m.id=a.media_id) INNER JOIN adx.sizes AS s ON (s.id=a.size_id) WHERE a.media_id=%d AND a.user_id=%d ORDER BY a.id DESC`, req.MediaId, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	var adzones []*common.Adzone
	for _, row := range rows {
		placeholderImg := row.Str(11)
		placeholderUrl := row.Str(12)
		var placeholder *common.PrivateAuctionCreative
		if placeholderImg != "" && placeholderUrl != "" {
			placeholder = &common.PrivateAuctionCreative{
				Url: placeholderUrl,
				Img: placeholderImg,
			}
			placeholder.ImgUrl = placeholder.GetImgUrl(Config)
		}
		adzone := &common.Adzone{
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
			Placeholder:  placeholder,
			Media: common.Media{
				Id:           row.Uint64(13),
				Title:        row.Str(14),
				Domain:       row.Str(15),
				OnlineStatus: row.Uint(16),
			},
			InsertedAt: row.ForceLocaltime(17),
			UpdatedAt:  row.ForceLocaltime(18),
		}
		adzone.EmbedCode = adzone.GetEmbedCode(Config)
		adzones = append(adzones, adzone)
	}
	if len(adzones) > 0 {
		rows, _, err = db.Query(`SELECT
	pa.adzone_id ,
	COUNT(*) AS auctions
FROM
	adx.private_auctions AS pa
INNER JOIN adx.adzones AS a ON ( a.id = pa.adzone_id )
INNER JOIN adx.medias AS m ON (m.id = a.media_id)
WHERE
	a.user_id = %d
AND a.online_status = 1
AND m.online_status = 1
AND pa.audit_status = 0
GROUP BY
	pa.adzone_id`, user.Id)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		auctionsMap := make(map[uint64]uint, len(rows))
		for _, row := range rows {
			auctionsMap[row.Uint64(0)] = row.Uint(1)
		}
		for _, adzone := range adzones {
			if auctions, found := auctionsMap[adzone.Id]; found {
				adzone.UnverifiedAuctions = auctions
			}
		}
	}
	c.JSON(http.StatusOK, adzones)
	return

}
