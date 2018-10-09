package promotion

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/tools/shorturl"
	"github.com/tokenme/adx/utils"
	"net/http"
	"time"
)

type AddRequest struct {
	AdzoneId  uint64 `form:"adzone_id" json:"adzone_id" binding:"required"`
	AirdropId uint64 `form:"airdrop_id" json:"airdrop_id" binding:"required"`
}

func AddHandler(c *gin.Context) {
	var req AddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT channel_id FROM adx.airdrop_adzones WHERE id=%d AND user_id=%d LIMIT 1`, req.AdzoneId, user.Id)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "invalid adzone", c) {
		return
	}
	channelId := rows[0].Uint64(0)
	rows, _, err = db.Query(`SELECT id, airdrop_id, adzone_id FROM adx.promotions WHERE user_id=%d AND airdrop_id=%d AND adzone_id=%d LIMIT 1`, user.Id, req.AirdropId, req.AdzoneId)
	if CheckErr(err, c) {
		return
	}
	var promotionId uint64
	if len(rows) > 0 {
		promotionId = rows[0].Uint64(0)
	} else {
		_, ret, err := db.Query(`INSERT INTO adx.promotions (user_id, airdrop_id, adzone_id, channel_id) VALUES (%d, %d, %d, %d)`, user.Id, req.AirdropId, req.AdzoneId, channelId)
		if CheckErr(err, c) {
			return
		}
		promotionId = ret.InsertId()
	}

	query := `SELECT
			uw.wallet,
			uw.salt
		FROM adx.user_wallets AS uw
		WHERE uw.user_id=%d AND uw.is_main=1`
	rows, _, err = db.Query(query, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	row := rows[0]
	wallet := row.Str(0)
	salt := row.Str(1)
	privateKey, _ := utils.AddressDecrypt(wallet, salt, Config.TokenSalt)
	publicKey, _ := eth.AddressFromHexPrivateKey(privateKey)

	promo := common.PromotionProto{
		Id:        promotionId,
		UserId:    user.Id,
		AirdropId: req.AirdropId,
		AdzoneId:  req.AdzoneId,
		ChannelId: channelId,
		Referrer:  publicKey,
	}

	promoKey, err := common.EncodePromotion([]byte(Config.LinkSalt), promo)
	if CheckErr(err, c) {
		return
	}
	link := fmt.Sprintf("%s/promo/%s", Config.BaseUrl, promoKey)
	shortURL, err := shorturl.Sina(link)
	if err == nil && shortURL != "" {
		link = shortURL
	}
	promotion := common.Promotion{
		Id:        promotionId,
		UserId:    user.Id,
		Airdrop:   &common.Airdrop{Id: req.AirdropId},
		AdzoneId:  req.AdzoneId,
		ChannelId: channelId,
		Link:      link,
		Key:       promoKey,
		Inserted:  time.Now(),
	}
	c.JSON(http.StatusOK, promotion)
}
