package promotion

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"time"
)

func StatsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	promotionId, err := Uint64NonZero(c.Query("promotion_id"), "missing promotion id")
	if CheckErr(err, c) {
		return
	}
	var stats []common.PromotionStats
	db := Service.Db

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	var (
		startDate time.Time
		endDate   time.Time
		today     = utils.TimeToDate(time.Now())
	)

	if endDateStr == "" {
		endDate = today
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			endDate = today
		}
	}

	if startDateStr == "" {
		startDate = endDate.AddDate(0, 0, -30)
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			startDate = endDate.AddDate(0, 0, -30)
		}
	}

	if startDate.After(endDate) {
		startDate = endDate.AddDate(0, 0, -30)
	}
	rows, _, err := db.Query(`SELECT SUM(ps.pv), SUM(ps.submissions), SUM(ps.transactions), SUM(ps.give_out), SUM(ps.bonus), SUM(ps.commission_fee), t.decimals, ps.record_on FROM adx.promotion_stats AS ps INNER JOIN adx.airdrops AS a ON (a.id=ps.airdrop_id) INNER JOIN adx.tokens AS t ON (t.address=a.token_address) WHERE ps.promotion_id=%d AND ps.record_on>='%s' AND ps.record_on<='%s' AND ps.promoter_id=%d GROUP BY ps.record_on ORDER BY ps.record_on ASC`, promotionId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), user.Id)
	if CheckErr(err, c) {
		return
	}
	var summary common.PromotionStats
	for _, row := range rows {
		s := common.PromotionStats{
			Pv:            row.Uint64(0),
			Submissions:   row.Uint64(1),
			Transactions:  row.Uint64(2),
			GiveOut:       row.Uint64(3),
			Bonus:         row.Uint64(4),
			CommissionFee: row.Uint64(5),
			Decimals:      row.Uint(6),
			RecordOn:      row.ForceLocaltime(7),
		}
		stats = append(stats, s)
		summary.Pv += s.Pv
		summary.Submissions += s.Submissions
		summary.Transactions += s.Transactions
		summary.GiveOut += s.GiveOut
		summary.Bonus += s.Bonus
		summary.CommissionFee += s.CommissionFee
	}
	c.JSON(http.StatusOK, common.PromotionStatsWithSummary{Summary: summary, Stats: stats})
}
