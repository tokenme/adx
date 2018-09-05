package adzone

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/ziutek/mymysql/mysql"
	"net/http"
	"strings"
	"time"
)

type ListRequest struct {
	MediaId uint64 `form:"media_id" json:"media_id" binding:"required"`
}

type Response struct {
	Id        int     `json:"id"`
	MediaName string  `json:"media_name"`
	Count     int     `json:"count"`
	LowPrice  string  `json:"low_price"`
	HighPrice string  `json:"high_price"`
	PV        uint64  `json:"pv"`
	UV        uint64  `json:"uv"`
	Clicks    uint64  `json:"clicks"`
	Ctr       float64 `json:"ctr"`
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

	if Check(user.IsPublisher != 1 && user.IsAdmin != 1, "unauthorized", c) {
		return
	}

	db := Service.Db
	rows := []mysql.Row{}
	var err error
	if user.IsAdmin == 1 {
		rows, _, err = db.Query(`SELECT a.id, a.url, a.size_id, s.width, s.height, a.min_cpm, a.min_cpt, a.settlement, a.rolling, a.intro, a.online_status, a.placeholder_img, a.placeholder_url, m.id, m.title, m.domain, m.online_status, a.inserted_at, a.updated_at FROM adx.adzones AS a INNER JOIN adx.medias AS m ON (m.id=a.media_id) INNER JOIN adx.sizes AS s ON (s.id=a.size_id) WHERE a.media_id=%d ORDER BY a.id DESC`, req.MediaId)
	} else {
		rows, _, err = db.Query(`SELECT a.id, a.url, a.size_id, s.width, s.height, a.min_cpm, a.min_cpt, a.settlement, a.rolling, a.intro, a.online_status, a.placeholder_img, a.placeholder_url, m.id, m.title, m.domain, m.online_status, a.inserted_at, a.updated_at FROM adx.adzones AS a INNER JOIN adx.medias AS m ON (m.id=a.media_id) INNER JOIN adx.sizes AS s ON (s.id=a.size_id) WHERE a.media_id=%d AND a.user_id=%d ORDER BY a.id DESC`, req.MediaId, user.Id)
	}
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

func TrafficListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsPublisher != 1 && user.IsAdvertiser != 1 && user.IsAdmin != 1, "unauthorized", c) {
		return
	}
	db := Service.Db
	ch := Service.Clickhouse
	Query := `SELECT a.id,a.title,MAX(b.min_cpt),MIN(b.min_cpt) From adx.medias AS a INNER JOIN adx.adzones AS b ON(a.id = b.media_id) GROUP BY b.media_id`
	rows, Result, err := db.Query(Query)
	if CheckErr(err, c) {
		return
	}
	var (
		Now          = time.Now()
		endDateStr   = time.Now().Format("2006-01-02")
		startDateStr = Now.AddDate(0, 0, -7).Format("2006-01-02")
	)
	List := []Response{}
	for _, value := range rows {
		Response := Response{}
		Response.Id = value.Int(Result.Map(`id`))
		Response.MediaName = value.Str(Result.Map(`title`))
		Response.HighPrice = value.Str(Result.Map(`MAX(b.min_cpt)`))
		Response.LowPrice = value.Str(Result.Map(`MIN(b.min_cpt)`))
		Query = `SELECT COUNT(*) FROM adx.adzones WHERE media_id=%d`
		rows, _, err = db.Query(Query, Response.Id)
		if CheckErr(err, c) {
			return
		}
		Response.Count = rows[0].Int(0)
		Query = `SELECT LogDate, pv, uv, clicks
FROM 
(
    SELECT 
        LogDate, 
        COUNTDistinct(ReqId) AS pv, 
        COUNTDistinct(Cookie) AS uv 
    FROM adx.reqs 
    WHERE %s
    GROUP BY LogDate
) ANY LEFT JOIN (
    SELECT 
        LogDate, 
        COUNTDistinct(ReqId) AS clicks 
    FROM adx.clicks 
    WHERE %s
    GROUP BY LogDate
) USING LogDate
ORDER BY LogDate ASC;`
		var (
			wheres []string
			date   time.Time
			pv     uint64
			uv     uint64
			clicks uint64
		)
		wheres = append(wheres, fmt.Sprintf("MediaId=%d", value.Int(Result.Map(`id`))))
		wheres = append(wheres, fmt.Sprintf("LogDate>='%s' AND LogDate <='%s'", startDateStr, endDateStr))
		where := strings.Join(wheres, " AND ")
		rows, err := ch.Query(fmt.Sprintf(Query, where, where))
		fmt.Println(fmt.Sprintf(Query, where, where))
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		for rows.Next() {
			err := rows.Scan(&date, &pv, &uv, &clicks)
			if CheckErr(err, c) {
				raven.CaptureError(err, nil)
				return
			}
			Response.PV += pv
			Response.UV += uv
			Response.Clicks += clicks
		}
		if Response.Clicks != 0 && Response.PV != 0 {
			Response.Ctr = float64(Response.Clicks) / float64(Response.PV)
		} else {
			Response.Ctr = 0.0
		}
		List = append(List, Response)
	}
	c.JSON(http.StatusOK, gin.H{
		`data`: List,
	})
	return
}
