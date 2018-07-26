package stats

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"strings"
	"time"
)

type BrowserRequest struct {
	MediaId   uint64   `form:"media_id" json:"media_id"`
	AdzoneId  uint64   `form:"adzone_id" json:"adzone_id"`
	AuctionId uint64   `form:"auction_id" json:"auction_id"`
	DateRange []string `form:"dateRange" json:"dateRange"`
}

func BrowserHandler(c *gin.Context) {
	var req BrowserRequest
	opt := c.Query("options")
	json.Unmarshal([]byte(opt), &req)
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	ch := Service.Clickhouse
	if Check(ch == nil, "stats server down", c) {
		return
	}

	var (
		startDate time.Time
		endDate   time.Time
		err       error
		wheres    []string
	)
	if len(req.DateRange) == 2 {
		startDate, err = time.Parse("2006-01-02", req.DateRange[0])
		if err != nil {
			endDate = utils.TimeToDate(time.Now())
			startDate = endDate.AddDate(0, -1, 0)
		} else {
			endDate, err = time.Parse("2006-01-02", req.DateRange[1])
		}
		if err != nil || endDate.Before(startDate) || endDate.After(startDate.AddDate(0, 3, 0)) {
			endDate = utils.TimeToDate(time.Now())
			startDate = endDate.AddDate(0, -1, 0)
		}
	} else {
		endDate = utils.TimeToDate(time.Now())
		startDate = endDate.AddDate(0, -1, 0)
	}
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	if Check(user.IsPublisher != 1 && user.IsAdvertiser != 1, "unauthorized", c) {
		return
	}
	if user.IsPublisher == 1 {
		wheres = append(wheres, fmt.Sprintf("PublisherId=%d", user.Id))
	} else if user.IsAdvertiser == 1 {
		wheres = append(wheres, fmt.Sprintf("AdvertiserId=%d", user.Id))
	}

	if req.MediaId > 0 {
		wheres = append(wheres, fmt.Sprintf("MediaId=%d", req.MediaId))
	}

	if req.AdzoneId > 0 {
		wheres = append(wheres, fmt.Sprintf("AdzoneId=%d", req.AdzoneId))
	}

	if req.AuctionId > 0 {
		wheres = append(wheres, fmt.Sprintf("AuctionId=%d", req.AuctionId))
	}

	if startDateStr == endDateStr {
		wheres = append(wheres, fmt.Sprintf("LogDate='%s'", startDateStr))
	} else {
		wheres = append(wheres, fmt.Sprintf("LogDate>='%s' AND LogDate <='%s'", startDateStr, endDateStr))
	}
	query := `SELECT BrowserName, pv, uv, clicks
FROM 
(
    SELECT 
        BrowserName, 
        COUNTDistinct(ReqId) AS pv, 
        COUNTDistinct(Cookie) AS uv 
    FROM adx.reqs 
    WHERE %s
    GROUP BY BrowserName
) ANY LEFT JOIN (
    SELECT 
        BrowserName, 
        COUNTDistinct(ReqId) AS clicks 
    FROM adx.clicks 
    WHERE %s
    GROUP BY BrowserName
) USING BrowserName
ORDER BY pv ASC LIMIT 10;`
	where := strings.Join(wheres, " AND ")
	rows, err := ch.Query(fmt.Sprintf(query, where, where))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	var osStats []common.OsStats
	for rows.Next() {
		var (
			name   string
			pv     uint64
			uv     uint64
			clicks uint64
			ctr    float64
		)
		err := rows.Scan(&name, &pv, &uv, &clicks)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		if pv > 0 {
			ctr = float64(clicks) / float64(pv)
		}
		osStats = append(osStats, common.OsStats{
			Name: name,
			Stats: common.Stats{
				Pv:     pv,
				Uv:     uv,
				Clicks: clicks,
				Ctr:    ctr,
			},
		})
	}
	c.JSON(http.StatusOK, osStats)
}
