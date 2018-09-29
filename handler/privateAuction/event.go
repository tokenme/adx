package privateAuction

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"strings"
)

type EventRequest struct {
	AdzoneId uint64 `form:"adzone_id" json:"adzone_id"`
	MediaId  uint64 `form:"media_id" json:"media_id"`
}

func EventHandler(c *gin.Context) {
	var req EventRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if Check(user.IsPublisher != 1 && user.IsAdmin != 1, "unauthorized", c) {
		return
	}

	db := Service.Db
	query := `SELECT
	pa.id AS id,
	pa.audit_status AS audit_status,
	pa.reject_reason AS reject_reason,
	pa.start_on AS start_on,
	pa.end_on AS end_on,
	pa.inserted_at AS inserted_at,
	pa.updated_at AS updated_at,
	pa.title AS title,
	pa.price AS price,
	pa.online_status AS online_status,
	pa.media_id AS media_id,
	m.title AS media_title,
	m.domain AS media_domain,
	pa.adzone_id AS adzone_id,
	s.id AS size_id,
	s.width AS size_width,
	s.height AS size_height
FROM
	adx.private_auctions AS pa 
INNER JOIN adx.adzones AS a ON (a.id = pa.adzone_id)
INNER JOIN adx.medias AS m ON (m.id = pa.media_id)
INNER JOIN adx.sizes AS s ON (s.id = a.size_id)
WHERE pa.online_status=1 AND pa.audit_status IN (0, 1) AND pa.end_on>=DATE(NOW()) AND a.user_id=%d%s
ORDER BY pa.start_on ASC, pa.end_on DESC, pa.price DESC
`

	var (
		where  string
		wheres []string
	)
	if req.AdzoneId > 0 {
		wheres = append(wheres, fmt.Sprintf("pa.adzone_id=%d", req.AdzoneId))
	} else if req.MediaId > 0 {
		wheres = append(wheres, fmt.Sprintf("pa.media_id=%d", req.MediaId))
	}
	if len(wheres) > 0 {
		where = fmt.Sprintf(" AND %s ", strings.Join(wheres, " AND "))
	}
	rows, _, err := db.Query(query, user.Id, where)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	var auctions []*common.PrivateAuction
	for _, row := range rows {
		auction := &common.PrivateAuction{
			Id:          row.Uint64(0),
			AuditStatus: row.Uint(1),
			StartTime:   row.ForceLocaltime(3),
			EndTime:     row.ForceLocaltime(4),
			Title:       row.Str(7),
			Price:       row.ForceFloat(8),
			Adzone: common.Adzone{
				Id: row.Uint64(13),
				Size: common.Size{
					Id:     row.Uint(14),
					Width:  row.Uint(15),
					Height: row.Uint(16),
				},
				Media: common.Media{
					Id:     row.Uint64(10),
					Title:  row.Str(11),
					Domain: row.Str(12),
				},
			},
		}
		auction.Cost = auction.GetCost()
		auctions = append(auctions, auction)
	}
	c.JSON(http.StatusOK, auctions)
}
