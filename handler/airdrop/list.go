package airdrop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net/http"
	"strings"
	"sync"
	"time"
)

const DEFAULT_PAGE_SIZE uint64 = 10

func ListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	var (
		where  string
		wheres []string
	)
	if user.IsPublisher != 0 || user.IsAdmin != 0 {
		var subWhere []string
		if user.IsAdmin == 0 {
			subWhere = append(subWhere, fmt.Sprintf("a.user_id=%d", user.Id))
		}
		status, _ := Uint64Value(c.Query("status"), 10)
		if status <= 2 {
			subWhere = append(subWhere, fmt.Sprintf("a.status=%d", status))
		}
		balanceStatus, _ := Uint64Value(c.Query("balance_status"), 10)
		if balanceStatus <= 4 {
			subWhere = append(subWhere, fmt.Sprintf("a.balance_status=%d", balanceStatus))
		}
		if len(subWhere) > 0 {
			wheres = append(wheres, strings.Join(subWhere, " AND "))
		}
	} else {
		now := time.Now().Format("2006-01-02")
		wheres = append(wheres, fmt.Sprintf("a.status=1 AND a.balance_status=0 AND a.start_date<='%s' AND end_date >='%s'", now, now))
	}
	if len(wheres) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(wheres, " OR "))
	}
	page, _ := Uint64Value(c.Query("page"), 1)
	if page == 0 {
		page = 1
	}
	pageSize, _ := Uint64Value(c.Query("page_size"), DEFAULT_PAGE_SIZE)
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	offset := (page - 1) * pageSize

	db := Service.Db
	rows, _, err := db.Query(`SELECT a.id, a.user_id, a.title, a.wallet, a.salt, t.address, t.name, t.symbol, t.decimals, t.protocol, a.gas_price, a.gas_limit, a.commission_fee, a.give_out, a.bonus, a.status, a.balance_status, a.start_date, a.end_date, a.drop_date, a.telegram_group, a.require_email, a.max_submissions, a.no_drop, a.reply_msg, a.inserted, a.updated, a.wallet_val_t, a.wallet_rule FROM adx.airdrops AS a INNER JOIN adx.tokens AS t ON (t.address=a.token_address) %s ORDER BY a.id DESC LIMIT %d, %d`, where, offset, pageSize)
	if CheckErr(err, c) {
		return
	}
	var airdrops []*common.Airdrop
	var wg sync.WaitGroup
	for _, row := range rows {
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
			wg.Add(1)
			go func(airdrop *common.Airdrop, c *gin.Context) {
				defer wg.Done()
				airdrop.CheckBalance(Service.Geth, c)
			}(airdrop, c)
		}
		airdrops = append(airdrops, airdrop)
	}
	wg.Wait()
	var val []string
	for _, a := range airdrops {
		if a.Token.Protocol == "ERC20" && a.NoDrop == 0 {
			val = append(val, fmt.Sprintf("(%d, %d)", a.Id, a.BalanceStatus))
		}
	}
	if len(val) > 0 {
		_, _, err = db.Query(`INSERT INTO adx.airdrops (id, balance_status) VALUES %s ON DUPLICATE KEY UPDATE balance_status=VALUES(balance_status)`, strings.Join(val, ","))
		if CheckErr(err, c) {
			return
		}
	}
	c.JSON(http.StatusOK, airdrops)
}
