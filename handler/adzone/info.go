package adzone

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"log"
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
	if Check(user.IsPublisher !=1 && user.IsAdvertiser !=1,"You don't have permission to access",c) {
		return
	}
	db := Service.Db
	query := `SELECT
	a.id ,
	a.url ,
	a.size_id ,
	s.width ,
	s.height ,
	a.min_cpm ,
	a.min_cpt ,
	a.settlement ,
	a.rolling ,
	a.intro ,
	a.online_status ,
	a.placeholder_url ,
	a.placeholder_img ,
	m.id ,
	m.title ,
	m.domain ,
	m.online_status ,
	a.inserted_at ,
	a.updated_at ,
	a.advantage,
	a.location,
	a.traffic
FROM
	adx.adzones AS a
INNER JOIN adx.medias AS m ON ( m.id = a.media_id )
INNER JOIN adx.sizes AS s ON ( s.id = a.size_id )
WHERE a.id=%d`
	rows, _, err := db.Query(query, req.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	row := rows[0]
	var placeholder *common.PrivateAuctionCreative
	placeholderUrl := row.Str(11)
	placeholderImg := row.Str(12)
	if placeholderUrl != "" && placeholderImg != "" {
		placeholder = &common.PrivateAuctionCreative{
			Url: placeholderUrl,
			Img: placeholderImg,
		}
		placeholder.ImgUrl = placeholder.GetImgUrl(Config)
	}
	adzone := common.Adzone{
		Id:  row.Uint64(0),
		Url: row.Str(1),
		Size: common.Size{
			Id:     row.Uint(2),
			Width:  row.Uint(3),
			Height: row.Uint(4),
		},
		SuggestCPT:   row.ForceFloat(6),
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
		InsertedAt:  row.ForceLocaltime(17),
		UpdatedAt:   row.ForceLocaltime(18),
		Advantage:row.Str(19),
		Location:row.Str(20),
		Traffic:row.Str(21),
	}
	unavailableDays, err := adzone.GetUnavailableDays(Service)
	if CheckErr(err, c) {
		return
	}
	adzone.UnavailableDays = unavailableDays
	query = `SELECT
	pa.adzone_id ,
	MAX(pa.price) AS price
FROM
	adx.private_auctions AS pa
WHERE
	pa.audit_status = 0
AND pa.online_status = 1
AND pa.adzone_id = %d
GROUP BY
	pa.adzone_id`
	rows, _, err = db.Query(query, adzone.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if len(rows) > 0 {
		row := rows[0]
		adzone.SuggestCPT = row.ForceFloat(1) * Config.AuctionRate / 100
	}

	c.JSON(http.StatusOK, adzone)
	return

}
