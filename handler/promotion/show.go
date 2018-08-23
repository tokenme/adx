package promotion

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"github.com/tokenme/adx/utils/token"
	"net/http"
	"time"
)

type ShowResponse struct {
	Airdrop *common.Airdrop `json:"airdrop,omitempty"`
	Code    token.Token     `json:"verify_code,omitempty"`
}

func ShowHandler(c *gin.Context) {
	cryptedKey := c.Param("key")
	if Check(cryptedKey == "", "missing key", c) {
		return
	}
	proto, err := common.DecodePromotion([]byte(Config.LinkSalt), cryptedKey)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT a.id, a.user_id, a.title, a.wallet, a.salt, t.address, t.name, t.symbol, t.decimals, t.protocol, t.client_ios, t.client_android, t.website, a.gas_price, a.gas_limit, a.commission_fee, a.give_out, a.bonus, a.status, a.balance_status, a.start_date, a.end_date, a.telegram_group, a.require_email, a.max_submissions, a.no_drop, a.inserted, a.updated, a.intro, a.promotion_page FROM adx.airdrops AS a INNER JOIN adx.tokens AS t ON (t.address=a.token_address) INNER JOIN adx.promotions AS p ON (p.airdrop_id=a.id) WHERE a.id=%d AND p.id=%d AND p.user_id=%d AND p.adzone_id=%d AND p.channel_id=%d`, proto.AirdropId, proto.Id, proto.UserId, proto.AdzoneId, proto.ChannelId)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "missing airdrop", c) {
		return
	}
	row := rows[0]
	wallet := row.Str(3)
	salt := row.Str(4)
	privateKey, _ := utils.AddressDecrypt(wallet, salt, Config.TokenSalt)
	publicKey, _ := eth.AddressFromHexPrivateKey(privateKey)
	airdrop := &common.Airdrop{
		Id:            row.Uint64(0),
		User:          common.User{Id: row.Uint64(1)},
		Title:         row.Str(2),
		Wallet:        publicKey,
		WalletPrivKey: privateKey,
		Token: common.Token{
			Address:       row.Str(5),
			Name:          row.Str(6),
			Symbol:        row.Str(7),
			Decimals:      row.Uint(8),
			Protocol:      row.Str(9),
			ClientIOS:     row.Str(10),
			ClientAndroid: row.Str(11),
			Website:       row.Str(12),
		},
		GasPrice:       row.Uint64(13),
		GasLimit:       row.Uint64(14),
		CommissionFee:  row.Uint64(15),
		GiveOut:        row.Uint64(16),
		Bonus:          row.Uint(17),
		Status:         row.Uint(18),
		BalanceStatus:  row.Uint(19),
		StartDate:      row.ForceLocaltime(20),
		EndDate:        row.ForceLocaltime(21),
		TelegramGroup:  row.Str(22),
		RequireEmail:   row.Uint(23),
		MaxSubmissions: row.Uint(24),
		NoDrop:         row.Uint(25),
		Inserted:       row.ForceLocaltime(26),
		Updated:        row.ForceLocaltime(27),
		TelegramBot:    Config.TelegramBotName,
		Intro:          row.Str(28),
		PromotionPage:  row.Str(29),
	}
	today := utils.TimeToDate(time.Now())
	if airdrop.StartDate.After(today) {
		airdrop.Status = common.AirdropStatusNotStart
	}
	if airdrop.EndDate.Before(today) {
		airdrop.Status = common.AirdropStatusFinished
	}
	var verifyCode token.Token
	if airdrop.Status == common.AirdropStatusStart {
		if airdrop.Token.Protocol == "ERC20" && airdrop.NoDrop == 0 {
			airdrop.CheckBalance(Service.Geth, c)
			_, _, err = db.Query(`UPDATE adx.airdrops SET balance_status=%d WHERE id=%d`, airdrop.BalanceStatus, airdrop.Id)
			if CheckErr(err, c) {
				return
			}
		}
		var submissions uint
		if airdrop.MaxSubmissions > 0 {
			rows, _, _ := db.Query(`SELECT SUM(submissions) FROM adx.promotion_stats WHERE airdrop_id=%d`, airdrop.Id)
			if len(rows) > 0 {
				submissions = rows[0].Uint(0)
			}
		}
		if (airdrop.Token.Protocol != "ERC20" || airdrop.BalanceStatus == common.AirdropBalanceStatusOk || airdrop.NoDrop == 1) && airdrop.MaxSubmissions >= submissions {
			for {
				verifyCode = token.New()
				_, _, err := db.Query(`INSERT INTO adx.codes (id, promotion_id, adzone_id, channel_id, promoter_id, airdrop_id) VALUES (%d, %d, %d, %d, %d, %d)`, verifyCode, proto.Id, proto.AdzoneId, proto.ChannelId, proto.UserId, airdrop.Id)
				if err != nil {
					continue
				}
				break
			}
		}
	}
	Tracker.Promotion.Pv(proto)

	c.JSON(http.StatusOK, ShowResponse{Airdrop: airdrop, Code: verifyCode})
}
