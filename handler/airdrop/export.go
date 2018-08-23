package airdrop

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"strconv"
	"time"
)

func SubmissionExportHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	airdropId, err := Uint64NonZero(c.Query("airdrop_id"), "missing airdrop id")
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	var where string
	if user.IsAdmin != 1 {
		where = fmt.Sprintf(" AND user_id=%d", user.Id)
	}
	rows, _, err := db.Query(`SELECT a.start_date, a.end_date, require_email FROM adx.airdrops AS a INNER JOIN adx.tokens AS t ON (t.address=a.token_address) WHERE a.id=%d%s`, airdropId, where)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}
	airdropStartDate := rows[0].ForceLocaltime(0)
	airdropEndDate := rows[0].ForceLocaltime(1)
	requireEmail := rows[0].Uint(2)
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
	query := `SELECT
	wallet ,
	referrer ,
	status ,
	telegram_user_id ,
	telegram_username ,
	telegram_user_firstname ,
	telegram_user_lastname ,
	tx ,
	inserted ,
	updated
FROM
	adx.airdrop_submissions
WHERE
	airdrop_id = %d
UNION
SELECT wallet,
    referrer,
    -1,
    -1,
    '',
    '',
    '',
    '',
    inserted,
    updated
FROM adx.airdrop_wallets AS aw
WHERE airdrop_id = %d
  AND NOT EXISTS (SELECT 1 FROM adx.airdrop_submissions AS a_sub WHERE a_sub.airdrop_id = aw.airdrop_id)
ORDER BY inserted ASC`
	fields := []string{"wallet", "referrer", "status", "telegram_user_id", "telegram_username", "telegram_user_firstname", "telegram_user_lastname", "tx", "inserted", "updated"}
	if requireEmail > 0 {
		query = `SELECT
	wallet ,
	referrer ,
	status ,
	email,
	telegram_user_id ,
	telegram_username ,
	telegram_user_firstname ,
	telegram_user_lastname ,
	tx ,
	inserted ,
	updated
FROM
	adx.airdrop_submissions
WHERE
	airdrop_id = %d
UNION
SELECT wallet,
    referrer,
    -1,
    email,
    -1,
    '',
    '',
    '',
    '',
    inserted,
    updated
FROM adx.airdrop_wallets AS aw
WHERE airdrop_id = %d
  AND NOT EXISTS (SELECT 1 FROM adx.airdrop_submissions AS a_sub WHERE a_sub.airdrop_id = aw.airdrop_id)
ORDER BY inserted ASC`
		fields = []string{"wallet", "referrer", "status", "email", "telegram_user_id", "telegram_username", "telegram_user_firstname", "telegram_user_lastname", "tx", "inserted", "updated"}
	}
	rows, _, err = db.Query(query, airdropId, airdropId)
	if CheckErr(err, c) {
		return
	}
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	lines := [][]string{fields}
	for _, row := range rows {
		var status string
		switch row.Int(2) {
		case 0:
			status = "pending"
		case 1:
			status = "transfered"
		case 2:
			status = "success"
		case 3:
			status = "failed"
		case -1:
			status = "incomplete"
		}
		var line []string
		if requireEmail > 0 {
			line = []string{
				row.Str(0),
				row.Str(1),
				status,
				row.Str(3),
				strconv.FormatInt(row.Int64(4), 10),
				row.Str(5),
				row.Str(6),
				row.Str(7),
				row.Str(8),
				row.ForceLocaltime(9).Format(time.RFC3339),
				row.ForceLocaltime(10).Format(time.RFC3339),
			}
		} else {
			line = []string{
				row.Str(0),
				row.Str(1),
				status,
				strconv.FormatInt(row.Int64(3), 10),
				row.Str(4),
				row.Str(5),
				row.Str(6),
				row.Str(7),
				row.ForceLocaltime(8).Format(time.RFC3339),
				row.ForceLocaltime(9).Format(time.RFC3339),
			}
		}
		lines = append(lines, line)
	}
	w.WriteAll(lines)
	w.Flush()
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"airdrop-%d-%s-%ssubmisstions.csv\"", airdropId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	//c.Header("Content-Length", strconv.Itoa(buf.Len()))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}
