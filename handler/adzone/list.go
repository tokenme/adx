package adzone

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/ziutek/mymysql/mysql"
	"math/rand"
	"net/http"
	"time"
	"fmt"
)

type ListRequest struct {
	MediaId uint64 `form:"media_id" json:"media_id" binding:"required"`
}

type Response struct {
	Id        uint64 `json:"id"`
	MediaName string `json:"media_name"`
	PV        uint64 `json:"pv"`
	UV        uint64 `json:"uv"`
	Clicks    uint64 `json:"clicks"`
	Ctr       string `json:"ctr"`
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
		rows, _, err = db.Query(`SELECT a.id, a.url, a.size_id, s.width, s.height, a.min_cpm, a.min_cpt, a.settlement, a.rolling, a.intro, a.online_status, a.placeholder_img, a.placeholder_url, m.id, m.title, m.domain, m.online_status, a.inserted_at, a.updated_at,a.advantage, a.location, a.traffic FROM adx.adzones AS a INNER JOIN adx.medias AS m ON (m.id=a.media_id) INNER JOIN adx.sizes AS s ON (s.id=a.size_id) WHERE a.media_id=%d ORDER BY a.id DESC`, req.MediaId)
	} else {
		rows, _, err = db.Query(`SELECT a.id, a.url, a.size_id, s.width, s.height, a.min_cpm, a.min_cpt, a.settlement, a.rolling, a.intro, a.online_status, a.placeholder_img, a.placeholder_url, m.id, m.title, m.domain, m.online_status, a.inserted_at, a.updated_at, a.advantage, a.location, a.traffic FROM adx.adzones AS a INNER JOIN adx.medias AS m ON (m.id=a.media_id) INNER JOIN adx.sizes AS s ON (s.id=a.size_id) WHERE a.media_id=%d AND a.user_id=%d ORDER BY a.id DESC`, req.MediaId, user.Id)
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
			Intro:        row.Str(9),
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
			Advantage:  row.Str(19),
			Location:   row.Str(20),
			Traffic:    row.Str(21),
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

var Res []Response
var times int64

func TrafficListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsAdvertiser != 1, "unauthorized", c) {
		return
	}
	Nows := time.Now()
	if times < Nows.Unix() {
		Res = []Response{}
		times = Nows.AddDate(0, 0, 7).Unix()
		db := Service.Db
		//ch := Service.Clickhouse
		Query := `SELECT id,title From adx.medias WHERE online_status = 1  AND verified = 1`
		rows, Result, err := db.Query(Query)
		if CheckErr(err, c) {
			return
		}
		var (
		//Now= time.Now()
		//endDateStr= time.Now().Format("2006-01-02")
		//startDateStr= Now.AddDate(0, 0, -7).Format("2006-01-02")
		)
		rand.Seed(time.Now().Unix())
		a := []float64{2.00, 2.10, 2.20,
			2.30, 2.40, 2.50, 2.60, 2.70, 2.80, 2.90,
			3.00, 3.10, 3.20, 3.30, 3.40, 3.50, 3.60,
			3.70, 3.80, 3.90, 4.00}
		for _, value := range rows {
			Response := Response{}
			rands := rand.Int63n(int64(len(a[:])) - 1)
			Response.MediaName = value.Str(Result.Map(`title`))
			Response.Id = value.Uint64(Result.Map(`id`))
			/*
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
						wheres = append(wheres, fmt.Sprintf("MediaId=%d", Response.Id))
						wheres = append(wheres, fmt.Sprintf("LogDate>='%s' AND LogDate <='%s'", startDateStr, endDateStr))
						where := strings.Join(wheres, " AND ")

						rows, err := ch.Query(fmt.Sprintf(Query, where, where))
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
			*/
			Response.PV = uint64(rand.Int63n(50000) + int64(50000*a[rands]))
			Response.UV = uint64(float64(Response.PV) / a[rands])
			Response.Clicks = Response.PV / uint64(rand.Int63n(10)+10)
			Response.Ctr = fmt.Sprintf("%.2f", float64(Response.Clicks)/float64(Response.PV)*100) + "%"
			Res = append(Res, Response)

			/*
				if Response.Clicks != 0 && Response.PV != 0 {
					Response.Ctr = float64(Response.Clicks) / float64(Response.PV)
				} else {
					Response.Ctr = 0.0
				}
				List = append(List, Response)
			*/
		}
		sort(Res)
	}

	c.JSON(http.StatusOK, gin.H{
		`data`: Res,
	})
	return
}

func sort(Res []Response) {
	for i := 0; i < len(Res); i++ {
		for k := i; k < len(Res); k++ {
			if Res[i].PV < Res[k].PV {
				Res[i], Res[k] = Res[k], Res[i]
			}
		}
	}
}
