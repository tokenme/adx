package promotion

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
)

const DEFAULT_PAGE_SIZE uint64 = 10

func ListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
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
	query := `SELECT
	p.id ,
	p.adzone_id ,
	p.channel_id ,
	p.airdrop_id ,
	a.title ,
	a.wallet ,
	a.salt ,
	t.address ,
	t.name ,
	t.symbol ,
	t.decimals ,
	t.protocol ,
	a.give_out ,
	a.bonus, 
	a.status ,
	a.balance_status ,
	a.start_date ,
	a.end_date ,
	a.telegram_group ,
	c.name ,
	az.name
FROM
	adx.promotions AS p
INNER JOIN adx.airdrops AS a ON ( a.id = p.airdrop_id )
INNER JOIN adx.tokens AS t ON ( t.address = a.token_address )
INNER JOIN adx.channels AS c ON ( c.id = p.channel_id )
INNER JOIN adx.adzones AS az ON ( az.id = p.adzone_id )
WHERE
	p.user_id =%d
ORDER BY
	p.id DESC
LIMIT %d, %d`
	rows, _, err := db.Query(query, user.Id, offset, pageSize)
	if CheckErr(err, c) {
		return
	}
	var (
		promotions []common.Promotion
		airdropMap = make(map[uint64]*common.Airdrop)
	)
	var wg sync.WaitGroup
	for _, row := range rows {
		airdropId := row.Uint64(3)
		var airdrop *common.Airdrop
		if a, found := airdropMap[airdropId]; found {
			airdrop = a
		} else {
			wallet := row.Str(5)
			salt := row.Str(6)
			privateKey, _ := utils.AddressDecrypt(wallet, salt, Config.TokenSalt)
			publicKey, _ := eth.AddressFromHexPrivateKey(privateKey)
			airdrop = &common.Airdrop{
				Id:            airdropId,
				Title:         row.Str(4),
				Wallet:        publicKey,
				WalletPrivKey: privateKey,
				Token: common.Token{
					Address:  row.Str(7),
					Name:     row.Str(8),
					Symbol:   row.Str(9),
					Decimals: row.Uint(10),
					Protocol: row.Str(11),
				},
				GiveOut:       row.Uint64(12),
				Bonus:         row.Uint(13),
				Status:        row.Uint(14),
				BalanceStatus: row.Uint(15),
				StartDate:     row.ForceLocaltime(16),
				EndDate:       row.ForceLocaltime(17),
				TelegramGroup: row.Str(18),
			}
			airdropMap[airdropId] = airdrop
			if airdrop.Token.Protocol == "ERC20" {
				wg.Add(1)
				go func(airdrop *common.Airdrop, c *gin.Context) {
					defer wg.Done()
					airdrop.CheckBalance(Service.Geth, c)
				}(airdrop, c)
			}
		}
		promotion := common.Promotion{
			Id:          row.Uint64(0),
			AdzoneId:    row.Uint64(1),
			ChannelId:   row.Uint64(2),
			Airdrop:     airdrop,
			ChannelName: row.Str(19),
			AdzoneName:  row.Str(20),
		}
		promo := common.PromotionProto{
			Id:        promotion.Id,
			UserId:    user.Id,
			AirdropId: promotion.Airdrop.Id,
			AdzoneId:  promotion.AdzoneId,
			ChannelId: promotion.ChannelId,
		}
		promoKey, err := common.EncodePromotion([]byte(Config.LinkSalt), promo)
		if CheckErr(err, c) {
			return
		}
		promotion.Key = promoKey
		promotion.Link = fmt.Sprintf("%s/promo/%s", Config.BaseUrl, promoKey)
		promotions = append(promotions, promotion)
	}
	wg.Wait()
	var val []string
	for _, a := range airdropMap {
		if a.Token.Protocol == "ERC20" {
			val = append(val, fmt.Sprintf("(%d, %d)", a.Id, a.BalanceStatus))
		}
	}
	if len(val) > 0 {
		_, _, err = db.Query(`INSERT INTO adx.airdrops (id, balance_status) VALUES %s ON DUPLICATE KEY UPDATE balance_status=VALUES(balance_status)`, strings.Join(val, ","))
		if CheckErr(err, c) {
			return
		}
	}
	c.JSON(http.StatusOK, promotions)
}
