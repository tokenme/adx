package airdrop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"time"
)

func StatsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	airdropId, err := Uint64NonZero(c.Query("airdrop_id"), "missing airdrop id")
	if CheckErr(err, c) {
		return
	}
	var stats []common.AirdropStats
	db := Service.Db
	var where string
	if user.IsAdmin != 1 {
		where = fmt.Sprintf(" AND user_id=%d", user.Id)
	}
	rows, _, err := db.Query(`SELECT a.start_date, a.end_date, t.decimals FROM adx.airdrops AS a INNER JOIN adx.tokens AS t ON (t.address=a.token_address) WHERE a.id=%d%s`, airdropId, where)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, common.AirdropStatsWithSummary{Summary: common.AirdropStats{}, Stats: stats})
		return
	}
	airdropStartDate := rows[0].ForceLocaltime(0)
	airdropEndDate := rows[0].ForceLocaltime(1)
	tokenDecimals := rows[0].Uint(2)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var (
		startDate time.Time
		endDate   time.Time
	)
	if startDateStr == "" {
		startDate = airdropStartDate
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			startDate = airdropStartDate
		}
	}
	if endDateStr == "" {
		endDate = airdropEndDate
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			endDate = airdropEndDate
		}
	}
	if endDate.After(airdropEndDate) {
		endDate = airdropEndDate
	}
	if startDate.Before(airdropStartDate) {
		startDate = airdropStartDate
	}
	if startDate.After(endDate) {
		startDate = endDate.AddDate(0, 0, -30)
	}
	rows, _, err = db.Query(`SELECT SUM(pv), SUM(submissions), SUM(transactions), SUM(give_out), SUM(bonus), SUM(commission_fee), record_on FROM adx.promotion_stats WHERE airdrop_id=%d AND record_on>='%s' AND record_on<='%s' GROUP BY record_on ORDER BY record_on ASC`, airdropId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if CheckErr(err, c) {
		return
	}
	var summary common.AirdropStats
	for _, row := range rows {
		s := common.AirdropStats{
			Pv:            row.Uint64(0),
			Submissions:   row.Uint64(1),
			Transactions:  row.Uint64(2),
			GiveOut:       row.Uint64(3),
			Bonus:         row.Uint64(4),
			CommissionFee: row.Uint64(5),
			Decimals:      tokenDecimals,
			RecordOn:      row.ForceLocaltime(6),
		}
		stats = append(stats, s)
		summary.Pv += s.Pv
		summary.Submissions += s.Submissions
		summary.Transactions += s.Transactions
		summary.GiveOut += s.GiveOut
		summary.Bonus += s.Bonus
		summary.CommissionFee += s.CommissionFee
	}
	c.JSON(http.StatusOK, common.AirdropStatsWithSummary{Summary: summary, Stats: stats})
}
