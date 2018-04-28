package privateAuction

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"strconv"
	"strings"
)

type ListRequest struct {
	AdzoneId    uint64 `form:"adzone_id" json:"adzone_id"`
	MediaId     uint64 `form:"media_id" json:"media_id"`
	AuctionId   uint64 `form:"auction_id" json:"auction_id"`
	Page        uint   `form:"page" json:"page"`
	AuditStatus int    `form:"audit_status" json:"audit_status"`
}

type ListResponse struct {
	Total    uint                     `json:"total"`
	Page     uint                     `json:"page"`
	PageSize uint                     `json:"page_size"`
	Auctions []*common.PrivateAuction `json:"auctions"`
}

const MAX_PAGE_SIZE uint = 20

func ListHandler(c *gin.Context) {
	var req ListRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

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
WHERE pa.online_status=1%s
ORDER BY %s
LIMIT %d, %d
`
	countQuery := `SELECT 
	COUNT(1)
FROM
	adx.private_auctions AS pa
INNER JOIN adx.adzones AS a ON (a.id = pa.adzone_id)
WHERE pa.online_status=1%s`

	var (
		where   string
		wheres  []string
		page    = req.Page
		orderBy string
	)
	if user.IsAdvertiser == 1 {
		wheres = append(wheres, fmt.Sprintf("pa.user_id=%d", user.Id))
		orderBy = "pa.id DESC"
	} else if user.IsPublisher == 1 {
		wheres = append(wheres, fmt.Sprintf("a.user_id=%d", user.Id))
		orderBy = "pa.start_on ASC, pa.end_on DESC, pa.price DESC"
	}
	if req.AuctionId > 0 {
		wheres = append(wheres, fmt.Sprintf("pa.id=%d", req.AuctionId))
	} else if req.AdzoneId > 0 {
		wheres = append(wheres, fmt.Sprintf("pa.adzone_id=%d", req.AdzoneId))
	} else if req.MediaId > 0 {
		wheres = append(wheres, fmt.Sprintf("pa.media_id=%d", req.MediaId))
	}
	if req.AuditStatus >= 0 {
		wheres = append(wheres, fmt.Sprintf("pa.audit_status=%d", req.AuditStatus))
	}
	if len(wheres) > 0 {
		where = fmt.Sprintf(" AND %s ", strings.Join(wheres, " AND "))
	}
	if page == 0 {
		page = 1
	}
	limit := (page - 1) * MAX_PAGE_SIZE
	rows, _, err := db.Query(query, where, orderBy, limit, MAX_PAGE_SIZE)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	var auctionsMap = make(map[uint64]*common.PrivateAuction, len(rows))
	var auctionIds []string
	for _, row := range rows {
		auction := &common.PrivateAuction{
			Id:           row.Uint64(0),
			AuditStatus:  row.Uint(1),
			RejectReason: row.Str(2),
			StartTime:    row.ForceLocaltime(3),
			EndTime:      row.ForceLocaltime(4),
			InsertedAt:   row.ForceLocaltime(5),
			UpdatedAt:    row.ForceLocaltime(6),
			Title:        row.Str(7),
			Price:        row.ForceFloat(8),
			OnlineStatus: row.Uint(9),
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
		auctionsMap[auction.Id] = auction
		auctionIds = append(auctionIds, strconv.FormatUint(auction.Id, 10))
	}

	if len(auctionIds) > 0 {
		rows, _, err := db.Query(`SELECT id, auction_id, landing_page, img FROM adx.private_auction_creatives WHERE auction_id IN (%s)`, strings.Join(auctionIds, ","))
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		for _, row := range rows {
			creative := common.PrivateAuctionCreative{
				Id:        row.Uint64(0),
				AuctionId: row.Uint64(1),
				Url:       row.Str(2),
				Img:       row.Str(3),
			}
			creative.ImgUrl = creative.GetImgUrl(Config)
			if auction, found := auctionsMap[creative.AuctionId]; found {
				auction.Creatives = append(auction.Creatives, creative)
			}
		}
	}
	var auctions []*common.PrivateAuction
	for _, auction := range auctionsMap {
		auctions = append(auctions, auction)
	}

	rows, _, err = db.Query(countQuery, where)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	total := rows[0].Uint(0)
	resp := ListResponse{
		Total:    total,
		Auctions: auctions,
		Page:     page,
		PageSize: MAX_PAGE_SIZE,
	}
	c.JSON(http.StatusOK, resp)
}
