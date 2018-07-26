package airdrop

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type ShareRequest struct {
	Id     uint64 `form:"id" json:"id" binding:"required"`
	Wallet string `form:"wallet" json:"wallet"`
}

func ShareHandler(c *gin.Context) {
	var req ShareRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	/*
		userContext, exists := c.Get("USER")
		if Check(!exists && req.Wallet == "", "need login", c) {
			return
		}
		user := userContext.(common.User)
	*/
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
				a.updated,
				uw.user_id ,
				uw.salt,
				uw.wallet,
			FROM adx.airdrops AS a 
			INNER JOIN adx.user_wallets AS uw ON (uw.user_id=a.user_id)
			INNER JOIN adx.tokens AS t ON (t.address=a.token_address) 
			WHERE 
				a.id=%d 
				AND uw.token_type = 'ETH'
				AND uw.is_main = 1`
	rows, _, err := db.Query(query, req.Id)
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

	tokenOwner := common.User{
		Id:     row.Uint64(16),
		Salt:   row.Str(17),
		Wallet: row.Str(18),
	}
	ethBalance, err := tokenOwner.ETHBalance(c, Service, Config)
	if CheckErr(err, c) {
		return
	}
	tokenBalance, err := tokenOwner.TokenBalance(c, Service, Config, airdrop.Token.Address)
	if CheckErr(err, c) {
		return
	}
	if airdrop.TotalTokenNeeded().Cmp(tokenBalance) > 0 {

	}
	if airdrop.TotalETHNeeded(Config.Airdrop).Cmp(ethBalance) > 0 {

	}
	c.JSON(http.StatusOK, airdrop)
	return

}
