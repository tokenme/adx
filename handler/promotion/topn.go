package promotion

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type TopNList struct {
	Wallet string `json:"wallet" codec:"wallet"`
	Cnt    int    `json:"cnt" codec:"cnt"`
	Idx    int    `json:"idx" codec:"idx"`
}

func TopNHandler(c *gin.Context) {
	n, err := Uint64NonZero(c.Query("n"), "missing top n")
	if CheckErr(err, c) {
		return
	}
	if n > 500 {
		n = 500
	}

	airdropId, err := Uint64NonZero(c.Query("airdrop_id"), "missing airdrop id")
	if CheckErr(err, c) {
		return
	}

	db := Service.Db
	query := `SELECT IF(TRIM(IFNULL(asub.referrer, "")) = "", asub.wallet, asub.referrer) AS referrer_s, count(*) AS cnt 
	FROM adx.airdrop_submissions AS asub
  WHERE asub.airdrop_id = %d
	GROUP BY referrer_s
  ORDER BY cnt DESC, referrer_s
  LIMIT %d`
	rows, _, err := db.Query(query, airdropId, n)
	if CheckErr(err, c) {
		return
	}

	ret := []*TopNList{}
	for idx, row := range rows {
		ret = append(ret, &TopNList{
			Wallet: row.Str(0),
			Cnt:    row.Int(1),
			Idx:    idx + 1,
		})
	}

	summary := struct {
		TotalCnt        int `json:"total_cnt" codec:"total_cnt"`
		Submissions     int `json:"submissions" codec:"submissions"`
		SelfSubmissions int `json:"self_submissions" codec:"self_submissions"`
	}{}
	query = `SELECT 
	COUNT(1) AS total_cnt,
		SUM(IF(TRIM(IFNULL(asub.referrer, "")) != "" AND asub.wallet != asub.referrer, 1, 0)) AS submissions
		FROM adx.airdrop_submissions AS asub
		WHERE asub.airdrop_id = %d`
	rows, _, err = db.Query(query, airdropId)
	if CheckErr(err, c) {
		return
	}
	if len(rows) > 0 {
		summary.TotalCnt = rows[0].Int(0)
		summary.Submissions = rows[0].Int(1)
		summary.SelfSubmissions = summary.TotalCnt - summary.Submissions
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"airdrop_id": airdropId,
		"topn":       ret,
		"summary":    summary,
	})
}
