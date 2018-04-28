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

type BrowserTypeRequest struct {
	MediaId   uint64   `form:"media_id" json:"media_id"`
	AdzoneId  uint64   `form:"adzone_id" json:"adzone_id"`
	AuctionId uint64   `form:"auction_id" json:"auction_id"`
	DateRange []string `form:"dateRange" json:"dateRange"`
}

func BrowserTypeHandler(c *gin.Context) {
	var req BrowserTypeRequest
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

	ch := Service.Clickhouse
	if startDateStr == endDateStr {
		wheres = append(wheres, fmt.Sprintf("LogDate='%s'", startDateStr))
	} else {
		wheres = append(wheres, fmt.Sprintf("LogDate>='%s' AND LogDate <='%s'", startDateStr, endDateStr))
	}
	query := `SELECT BrowserType, pv, uv, clicks
FROM 
(
    SELECT 
        BrowserType, 
        COUNTDistinct(ReqId) AS pv, 
        COUNTDistinct(Cookie) AS uv 
    FROM adx.reqs 
    WHERE %s
    GROUP BY BrowserType
) ANY LEFT JOIN (
    SELECT 
        BrowserType, 
        COUNTDistinct(ReqId) AS clicks 
    FROM adx.clicks 
    WHERE %s
    GROUP BY BrowserType
) USING BrowserType
ORDER BY pv ASC;`
	where := strings.Join(wheres, " AND ")
	rows, err := ch.Query(fmt.Sprintf(query, where, where))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	var browserTypeStats []common.BrowserTypeStats
	for rows.Next() {
		var (
			browserType uint
			pv          uint64
			uv          uint64
			clicks      uint64
			ctr         float64
		)
		err := rows.Scan(&browserType, &pv, &uv, &clicks)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		if pv > 0 {
			ctr = float64(clicks) / float64(pv)
		}
		var name string
		switch browserType {
		case 0:
			name = "unknown"
		case 1:
			name = "desktop"
		case 2:
			name = "mobile"
		case 4:
			name = "console"
		}
		browserTypeStats = append(browserTypeStats, common.BrowserTypeStats{
			Id:   browserType,
			Name: name,
			Stats: common.Stats{
				Pv:     pv,
				Uv:     uv,
				Clicks: clicks,
				Ctr:    ctr,
			},
		})
	}
	c.JSON(http.StatusOK, browserTypeStats)
}
