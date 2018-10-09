package airdrop

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
)

func GetHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	airdropId, err := Uint64NonZero(c.Query("id"), "missing airdrop id")
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT a.id, a.user_id, a.title, a.wallet, a.salt, t.address, t.name, t.symbol, t.decimals, t.protocol, a.gas_price, a.gas_limit, a.commission_fee, a.give_out, a.bonus, a.status, a.balance_status, a.start_date, a.end_date, a.drop_date, a.telegram_group, a.require_email, a.max_submissions, a.no_drop, a.reply_msg, a.inserted, a.updated, a.wallet_val_t, a.wallet_rule FROM adx.airdrops AS a INNER JOIN adx.tokens AS t ON (t.address=a.token_address) WHERE a.id=%d`, airdropId)
	if CheckErr(err, c) {
		return
	}
	var airdrop *common.Airdrop
	if len(rows) > 0 {
		row := rows[0]
		wallet := row.Str(3)
		salt := row.Str(4)
		privateKey, _ := utils.AddressDecrypt(wallet, salt, Config.TokenSalt)
		publicKey, _ := eth.AddressFromHexPrivateKey(privateKey)
		airdrop = &common.Airdrop{
			Id:            row.Uint64(0),
			User:          common.User{Id: row.Uint64(1)},
			Title:         row.Str(2),
			Wallet:        publicKey,
			WalletPrivKey: privateKey,
			Token: common.Token{
				Address:  row.Str(5),
				Name:     row.Str(6),
				Symbol:   row.Str(7),
				Decimals: row.Uint(8),
				Protocol: row.Str(9),
			},
			GasPrice:       row.Uint64(10),
			GasLimit:       row.Uint64(11),
			CommissionFee:  row.Uint64(12),
			GiveOut:        row.Uint64(13),
			Bonus:          row.Uint(14),
			Status:         row.Uint(15),
			BalanceStatus:  row.Uint(16),
			StartDate:      row.ForceLocaltime(17),
			EndDate:        row.ForceLocaltime(18),
			DropDate:       row.ForceLocaltime(19),
			TelegramBot:    Config.TelegramBotName,
			TelegramGroup:  row.Str(20),
			RequireEmail:   row.Uint(21),
			MaxSubmissions: row.Uint(22),
			NoDrop:         row.Uint(23),
			ReplyMsg:       row.Str(24),
			Inserted:       row.ForceLocaltime(25),
			Updated:        row.ForceLocaltime(26),
			WalletValType:  uint8(row.Uint(27)),
			WalletRule:     row.Str(28),
		}
		if airdrop.Token.Protocol == "ERC20" {
			airdrop.CheckBalance(Service.Geth, c)
			if airdrop.NoDrop == 0 {
				_, _, err = db.Query(`UPDATE adx.airdrops SET balance_status=%d WHERE id=%d`, airdrop.BalanceStatus, airdrop.Id)
				if CheckErr(err, c) {
					return
				}
			}
		}
	}
	c.JSON(http.StatusOK, airdrop)
}
