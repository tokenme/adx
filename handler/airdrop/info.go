package airdrop

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type InfoRequest struct {
	Id uint64 `form:"id" json:"id" binding:"required"`
}

func InfoHandler(c *gin.Context) {
	var req InfoRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if Check(user.IsAdvertiser != 1, "unauthorized", c) {
		return
	}

	db := Service.Db

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
			WHERE a.id=%d AND a.user_id=%d`
	rows, _, err := db.Query(query, req.Id, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	if Check(len(rows) == 0, "not found", c) {
		return
	}
	row := rows[0]
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
	c.JSON(http.StatusOK, airdrop)
	return

}
