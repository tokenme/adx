package privateAuction

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"strings"
)

type AuditRequest struct {
	Id           uint64 `form:"id" json:"id" binding:"required"`
	AuditStatus  uint   `form:"audit_status" json:"audit_status" binding:"required"`
	RejectReason string `form:"reject_reason" json:"reject_reason"`
}

func AuditHandler(c *gin.Context) {
	var req AuditRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if Check(user.IsPublisher != 1, "unauthorized", c) {
		return
	}

	db := Service.Db
	if Check(req.AuditStatus != 1 && req.AuditStatus != 2, "invalid audit", c) {
		return
	}
	rows, _, err := db.Query(`SELECT pa.adzone_id, pa.start_on, pa.end_on FROM adx.private_auctions AS pa INNER JOIN adx.adzones AS a ON (a.id=pa.adzone_id) WHERE pa.id=%d AND a.user_id=%d LIMIT 1`, req.Id, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	adzoneId := rows[0].Uint64(0)
	startTime := rows[0].ForceLocaltime(1)
	endTime := rows[0].ForceLocaltime(2)
	startStr := startTime.Format("2006-01-02")
	endStr := endTime.Format("2006-01-02")
	if req.AuditStatus == 2 {
		rejectReason := utils.Normalize(req.RejectReason)
		if Check(rejectReason == "", "missing reject reason", c) {
			return
		}
		_, _, err := db.Query(`UPDATE adx.private_auctions SET audit_status=2, reject_reason='%s' WHERE id=%d AND audit_status=0`, db.Escape(rejectReason), req.Id)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
	} else {
		_, _, err := db.Query(`UPDATE adx.private_auctions SET audit_status=1, reject_reason=NULL WHERE id=%d AND audit_status=0`, req.Id)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		_, _, err = db.Query(`UPDATE adx.private_auctions SET audit_status=2, reject_reason='audit failed' WHERE adzone_id=%d AND audit_status=0 AND (start_on BETWEEN '%s' AND '%s' OR end_on BETWEEN '%s' AND '%s')`, adzoneId, startStr, endStr, startStr, endStr)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			_, _, err := db.Query(`UPDATE adx.private_auctions SET audit_status=0 WHERE id=%d`, req.Id)
			if CheckErr(err, c) {
				raven.CaptureError(err, nil)
			}
			return
		}
		day := startTime
		var val []string
		for {
			if day.After(endTime) {
				break
			}
			val = append(val, fmt.Sprintf("(%d, %d, '%s')", adzoneId, req.Id, db.Escape(day.Format("2006-01-02"))))
			day = day.AddDate(0, 0, 1)
		}
		if len(val) > 0 {
			_, _, err = db.Query(`INSERT IGNORE INTO adx.adzone_auction_days (adzone_id, auction_id, record_on) VALUES %s`, strings.Join(val, ","))
			if CheckErr(err, c) {
				raven.CaptureError(err, nil)
				_, _, err = db.Query(`UPDATE adx.private_auctions SET audit_status=0, reject_reason=NULL WHERE adzone_id=%d AND (start_on BETWEEN '%s' AND '%s' OR end_on BETWEEN '%s' AND '%s')`, adzoneId, startStr, endStr, startStr, endStr)
				if CheckErr(err, c) {
					raven.CaptureError(err, nil)
				}
				return
			}
		}
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
