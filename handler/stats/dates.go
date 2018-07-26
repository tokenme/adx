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

type DatesRequest struct {
	MediaId   uint64   `form:"media_id" json:"media_id"`
	AdzoneId  uint64   `form:"adzone_id" json:"adzone_id"`
	AuctionId uint64   `form:"auction_id" json:"auction_id"`
	DateRange []string `form:"dateRange" json:"dateRange"`
}

func DatesHandler(c *gin.Context) {
	var req DatesRequest
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

	var query string
	if startDateStr == endDateStr {
		wheres = append(wheres, fmt.Sprintf("LogDate='%s'", startDateStr))
		query = `SELECT h, pv, uv, clicks
FROM 
(
    SELECT 
        toHour(LogTime) AS h, 
        COUNTDistinct(ReqId) AS pv, 
        COUNTDistinct(Cookie) AS uv 
    FROM adx.reqs 
    WHERE %s
    GROUP BY h
) ANY LEFT JOIN (
    SELECT 
        toHour(LogTime) AS h, 
        COUNTDistinct(ReqId) AS clicks 
    FROM adx.clicks 
    WHERE %s
    GROUP BY h
) USING h
ORDER BY h ASC;`
	} else {
		wheres = append(wheres, fmt.Sprintf("LogDate>='%s' AND LogDate <='%s'", startDateStr, endDateStr))
		query = `SELECT LogDate, pv, uv, clicks
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
	}
	where := strings.Join(wheres, " AND ")
	rows, err := ch.Query(fmt.Sprintf(query, where, where))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	var (
		dateStatsMap = make(map[string]common.DateStats)
		hourStatsMap = make(map[uint]common.HourStats)
	)
	for rows.Next() {
		var (
			date   time.Time
			h      uint
			pv     uint64
			uv     uint64
			clicks uint64
			ctr    float64
		)
		if startDateStr == endDateStr {
			err := rows.Scan(&h, &pv, &uv, &clicks)
			if CheckErr(err, c) {
				raven.CaptureError(err, nil)
				return
			}
			if pv > 0 {
				ctr = float64(clicks) / float64(pv)
			}
			hourStatsMap[h] = common.HourStats{
				Hour: h,
				Stats: common.Stats{
					Pv:     pv,
					Uv:     uv,
					Clicks: clicks,
					Ctr:    ctr,
				},
			}
		} else {
			err := rows.Scan(&date, &pv, &uv, &clicks)
			if CheckErr(err, c) {
				raven.CaptureError(err, nil)
				return
			}
			if pv > 0 {
				ctr = float64(clicks) / float64(pv)
			}
			dateStr := date.Format("2006-01-02")
			dateStatsMap[dateStr] = common.DateStats{
				Date: dateStr,
				Stats: common.Stats{
					Pv:     pv,
					Uv:     uv,
					Clicks: clicks,
					Ctr:    ctr,
				},
			}
		}
	}
	if startDateStr == endDateStr {
		var (
			hourStats []common.HourStats
			h         uint
		)
		for h < 24 {
			if s, found := hourStatsMap[h]; found {
				hourStats = append(hourStats, s)
			} else {
				hourStats = append(hourStats, common.HourStats{
					Hour:  h,
					Stats: common.Stats{},
				})
			}
			h += 1
		}
		c.JSON(http.StatusOK, hourStats)
	} else {
		var (
			dateStats []common.DateStats
			date      = startDate
		)
		for !date.After(endDate) {
			dateStr := date.Format("2006-01-02")
			if s, found := dateStatsMap[dateStr]; found {
				dateStats = append(dateStats, s)
			} else {
				dateStats = append(dateStats, common.DateStats{
					Date:  dateStr,
					Stats: common.Stats{},
				})
			}
			date = date.AddDate(0, 0, 1)
		}
		c.JSON(http.StatusOK, dateStats)
	}
	return

}
