package airdrop

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

type ListRequest struct {
	Page         uint     `form:"page" json:"page"`
	OnlineStatus int      `form:"onlineStatus" json:"onlineStatus"`
	DateRange    []string `form:"dateRange" json:"dateRange"`
}

type ListResponse struct {
	Total    uint             `json:"total"`
	Page     uint             `json:"page"`
	PageSize uint             `json:"page_size"`
	Airdrops []common.Airdrop `json:"airdrops"`
}

const MAX_PAGE_SIZE uint = 20

func ListHandler(c *gin.Context) {
	var req ListRequest
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
	var (
		where  string
		wheres []string
		page   = req.Page
	)
	wheres = append(wheres, fmt.Sprintf("a.user_id=%d", user.Id))
	if req.OnlineStatus != 0 {
		wheres = append(wheres, fmt.Sprintf("a.online_status=%d", req.OnlineStatus))
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
	wheres = append(wheres, fmt.Sprintf("a.start_date>='%s' AND a.end_date<='%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	if len(wheres) > 0 {
		where = strings.Join(wheres, " AND ")
	}
	if page == 0 {
		page = 1
	}
	limit := (page - 1) * MAX_PAGE_SIZE

	query := `SELECT 
				a.id, 
				a.title, 
				t.address, 
				t.name, 
				t.symbol, 
				t.decimals,
				a.budget, 
				a.commission_fee, 
				a.give_out, 
				a.bonus, 
				a.online_status,
				a.start_date, 
				a.end_date, 
				a.telegram_group, 
				a.inserted, 
				a.updated 
			FROM adx.airdrops AS a 
			INNER JOIN adx.tokens AS t ON (t.address=a.token_address) 
			WHERE %s ORDER BY a.id DESC LIMIT %d, %d`

	countQuery := `SELECT 
	COUNT(1)
FROM
	adx.airdrops AS a
INNER JOIN adx.tokens AS t ON (t.address=a.token_address)
WHERE %s`
	rows, _, err := db.Query(query, where, limit, MAX_PAGE_SIZE)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	var airdrops []common.Airdrop
	for _, row := range rows {
		airdrop := common.Airdrop{
			Id:    row.Uint64(0),
			Title: row.Str(1),
			Token: common.Token{
				Address:  row.Str(2),
				Name:     row.Str(3),
				Symbol:   row.Str(4),
				Decimals: row.Uint(5),
			},
			Budget:        row.Uint64(6),
			CommissionFee: row.Uint64(7),
			GiveOut:       row.Uint64(8),
			Bonus:         row.Uint(9),
			OnlineStatus:  row.Int(10),
			StartDate:     row.ForceLocaltime(11),
			EndDate:       row.ForceLocaltime(12),
			TelegramGroup: row.Str(13),
			InsertedAt:    row.ForceLocaltime(14),
			UpdatedAt:     row.ForceLocaltime(15),
		}
		airdrops = append(airdrops, airdrop)
	}

	rows, _, err = db.Query(countQuery, where)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	total := rows[0].Uint(0)
	resp := ListResponse{
		Total:    total,
		Airdrops: airdrops,
		Page:     page,
		PageSize: MAX_PAGE_SIZE,
	}
	c.JSON(http.StatusOK, resp)
	return

}
