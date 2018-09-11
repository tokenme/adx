package adzone

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SearchRequest struct {
	SizeIds   []uint   `form:"sizes" json:"sizes"`
	Domain    string   `form:"domain" json:"domain"`
	DateRange []string `form:"dateRange" json:"dateRange"`
	MediaId   uint64   `form:"media_id" json:"media_id"`
}

func SearchHandler(c *gin.Context) {
	var req SearchRequest
	opt := c.Query("options")
	json.Unmarshal([]byte(opt), &req)
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if Check(user.IsAdvertiser != 1, "unauthorized", c) {
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
	m.id ,
	m.title ,
	m.domain ,
	m.online_status , 
	a.inserted_at ,
	a.updated_at
FROM
	adx.adzones AS a
INNER JOIN adx.medias AS m ON ( m.id = a.media_id )
INNER JOIN adx.sizes AS s ON ( s.id = a.size_id )
WHERE
	a.online_status = 1
AND m.online_status = 1
%s
ORDER BY
	a.id DESC`
	var (
		wheres    []string
		subWheres []string
	)
	if len(req.SizeIds) > 0 {
		var sizeIds []string
		for _, sizeId := range req.SizeIds {
			sizeIds = append(sizeIds, fmt.Sprintf("%d", sizeId))
		}
		wheres = append(wheres, fmt.Sprintf("a.size_id IN (%s)", strings.Join(sizeIds, ",")))
		subWheres = append(wheres, fmt.Sprintf("a2.size_id IN (%s)", strings.Join(sizeIds, ",")))
	}
	if req.MediaId > 0 {
		wheres = append(wheres, fmt.Sprintf("a.media_id=%d", req.MediaId))
		subWheres = append(wheres, fmt.Sprintf("a2.media_id=%d", req.MediaId))
	} else if req.Domain != "" {
		domain := req.Domain
		if strings.HasPrefix(domain, "http") {
			wheres = append(wheres, fmt.Sprintf("m.domain='%s'", db.Escape(domain)))
			subWheres = append(subWheres, fmt.Sprintf("m2.domain='%s'", db.Escape(domain)))
		} else {
			wheres = append(wheres, fmt.Sprintf("(m.domain='%s' OR m.domain='%s')", db.Escape(fmt.Sprintf("https://%s", domain)), db.Escape(fmt.Sprintf("http://%s", domain))))
			subWheres = append(subWheres, fmt.Sprintf("(m2.domain='%s' OR m2.domain='%s')", db.Escape(fmt.Sprintf("https://%s", domain)), db.Escape(fmt.Sprintf("http://%s", domain))))
		}
	}

	var (
		startDate time.Time
		endDate   time.Time
		err       error
	)
	if len(req.DateRange) == 2 {
		startDate, err = time.Parse("2006-01-02", req.DateRange[0])
		if err != nil {
			startDate = utils.TimeToDate(time.Now())
			endDate = startDate.AddDate(0, 2, 0)
		} else {
			endDate, err = time.Parse("2006-01-02", req.DateRange[1])
		}
		if err != nil || endDate.Before(startDate) || endDate.After(startDate.AddDate(0, 2, 0)) {
			startDate = utils.TimeToDate(time.Now())
			endDate = startDate.AddDate(0, 2, 0)
		}
	} else {
		startDate = utils.TimeToDate(time.Now())
		endDate = startDate.AddDate(0, 2, 0)
	}
	days := int(endDate.Sub(startDate).Hours())/24 + 1
	subWheres = append(subWheres, fmt.Sprintf("aad.record_on BETWEEN '%s' AND '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	var subWhere string
	if len(subWheres) > 0 {
		subWhere = strings.Join(subWheres, " AND ")
	}

	subQuery := fmt.Sprintf(`SELECT
	adzone_id ,
	COUNT(*) AS num
FROM
	( SELECT
		aad.adzone_id ,
		COUNT( aad.auction_id ) AS auctions ,
		a2.rolling AS rolling ,
		aad.record_on AS record_on
	FROM
		adx.adzone_auction_days AS aad
	INNER JOIN adx.adzones AS a2 ON ( a2.id = aad.adzone_id )
	INNER JOIN adx.medias AS m2 ON (m2.id = a2.media_id)
	WHERE
		%s
	GROUP BY
		aad.adzone_id ,
		aad.record_on
	HAVING
		auctions >= rolling ) AS tmp
	WHERE tmp.adzone_id = a.id
GROUP BY
	tmp.adzone_id
HAVING num >= %d`, subWhere, days)
	wheres = append(wheres, fmt.Sprintf("NOT EXISTS (%s)", subQuery))
	var where string
	if len(wheres) > 0 {
		where = fmt.Sprintf(" AND %s", strings.Join(wheres, " AND "))
	}
	rows, _, err := db.Query(query, where)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	var (
		adzones   []*common.Adzone
		adzoneIds []string
	)
	for _, row := range rows {
		adzone := &common.Adzone{
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
			Intro:         row.Str(9),
			OnlineStatus: row.Uint(10),
			Media: common.Media{
				Id:           row.Uint64(11),
				Title:        row.Str(12),
				Domain:       row.Str(13),
				OnlineStatus: row.Uint(14),
			},
			InsertedAt: row.ForceLocaltime(15),
			UpdatedAt:  row.ForceLocaltime(16),
		}
		adzoneIds = append(adzoneIds, strconv.FormatUint(adzone.Id, 10))
		adzones = append(adzones, adzone)
	}

	if len(adzoneIds) > 0 {
		adzoneUnavailableDays, err := common.AdzonesUnavailableDays(Service, adzoneIds)
		if CheckErr(err, c) {
			return
		}
		query = `SELECT
	pa.adzone_id ,
	MAX(pa.price) AS price
FROM
	adx.private_auctions AS pa
WHERE
	pa.audit_status = 0
AND pa.online_status = 1
AND pa.adzone_id IN (%s)
GROUP BY
	pa.adzone_id`
		rows, _, err := db.Query(query, strings.Join(adzoneIds, ","))
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		suggestPrices := make(map[uint64]float64, len(rows))
		for _, row := range rows {
			suggestPrices[row.Uint64(0)] = row.ForceFloat(1) * Config.AuctionRate / 100
		}

		for _, ad := range adzones {
			if price, found := suggestPrices[ad.Id]; found {
				ad.SuggestCPT = price
			}
			if days, found := adzoneUnavailableDays[ad.Id]; found {
				ad.UnavailableDays = days
			}
		}
	}

	c.JSON(http.StatusOK, adzones)
	return

}

type AdzoneMedia struct {
	Id          uint    `json:"id"`
	MediaName   string  `json:"media_name"`
	AdzoneCount uint64  `json:"adzone_count"`
	HighPrice   float64 `json:"high_price"`
	LowPrice    float64 `json:"low_price"`
}


func MediaListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsAdvertiser != 1, "unauthorized", c) {
		return
	}
	var req SearchRequest
	opt := c.Query("options")
	json.Unmarshal([]byte(opt), &req)
	db := Service.Db
	Where := []string{}
	Query:= ""
	if req.MediaId > 0 || req.Domain != "" || len(req.DateRange) > 0 || len(req.SizeIds) > 0 {
		if req.MediaId > 0 {
			Where = append(Where, fmt.Sprintf("media.id=%d", req.MediaId))
		}
		if len(req.SizeIds) > 0 {
			var sizeIds []string
			for _, sizeId := range req.SizeIds {
				sizeIds = append(sizeIds, fmt.Sprintf("%d", sizeId))
			}
			Where = append(Where, fmt.Sprintf("adzon.size_id IN (%s)", strings.Join(sizeIds, ",")))
		}
		if req.Domain != "" {
			Where = append(Where, fmt.Sprintf("adzon.url='%s'", req.Domain))
		}
		var (
			startDate time.Time
			endDate   time.Time
			err       error
		)
		if len(req.DateRange) == 2 {
			startDate, err = time.Parse("2006-01-02", req.DateRange[0])
			if CheckErr(err, c) {
				return
			}
			endDate, err = time.Parse("2006-01-02", req.DateRange[1])
			if CheckErr(err, c) {
				return
			}
		} else {
			startDate = utils.TimeToDate(time.Now())
			endDate = startDate.AddDate(0, 0, 7)
		}
		Where = append(Where, fmt.Sprintf("aaday.record_on BETWEEN '%s' AND '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
		WHERES := strings.Join(Where, " AND ")
		Querys := `SELECT media.id,media.title,MAX(adzon.min_cpt),MIN(adzon.min_cpt),COUNT(*)
FROM adx.medias AS media
INNER JOIN adx.adzones AS adzon ON ( adzon.media_id = media.id )
INNER JOIN adx.adzone_auction_days AS aaday ON ( aaday.adzone_id = adzon.id)
WHERE %s 
`
		Query = fmt.Sprintf(Querys,WHERES)
	} else {
		Query = `SELECT media.id,media.title,MAX(adzon.min_cpt),MIN(adzon.min_cpt) ,COUNT(*)
From adx.medias AS media INNER JOIN adx.adzones 
AS adzon ON(media.id = adzon.media_id) GROUP BY adzon.media_id`
	}
		rows, Result, err := db.Query(Query)
		if CheckErr(err, c) {
			return
		}
		AdzoneMediaList := []AdzoneMedia{}
		for _, value := range rows {
			Media := AdzoneMedia{}
			Media.Id = value.Uint(Result.Map(`id`))
			Media.MediaName = value.Str(Result.Map(`title`))
			Media.HighPrice = value.Float(Result.Map(`MAX(adzon.min_cpt)`))
			Media.LowPrice = value.Float(Result.Map(`MIN(adzon.min_cpt)`))
			Media.AdzoneCount = value.Uint64(Result.Map(`COUNT(*)`))
			AdzoneMediaList = append(AdzoneMediaList, Media)
		}

		c.JSON(http.StatusOK, gin.H{
			`Data`: AdzoneMediaList,
		})
	}


