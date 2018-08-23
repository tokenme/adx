package promotion

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/tools/shorturl"
	"github.com/tokenme/adx/utils/token"
	"net/http"
	"regexp"
	"strings"
)

var Email string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

type SubmitRequest struct {
	Wallet   string      `json:"wallet,omitempty"`
	Email    string      `json:"email,omitempty"`
	Code     token.Token `json:"verify_code,omitempty"`
	ProtoKey string      `json:"proto,omitempty"`
}

func SubmitHandler(c *gin.Context) {
	var req SubmitRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	proto, err := common.DecodePromotion([]byte(Config.LinkSalt), req.ProtoKey)
	if CheckErr(err, c) {
		return
	}

	db := Service.Db
	rows, _, err := db.Query(`SELECT t.protocol, a.require_email FROM adx.tokens AS t INNER JOIN adx.airdrops AS a ON (a.token_address=t.address) WHERE a.id=%d LIMIT 1`, proto.AirdropId)
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	protocol := rows[0].Str(0)
	requireEmail := rows[0].Uint(1)
	emailRegex := regexp.MustCompile(Email)
	if Check(requireEmail > 0 && (req.Email == "" || !emailRegex.MatchString(req.Email)), "invalid email address", c) {
		return
	}
	if Check(protocol == "ERC20" && (len(req.Wallet) != 42 || !strings.HasPrefix(req.Wallet, "0x")), "invalid wallet", c) {
		return
	}
	rows, _, err = db.Query("SELECT id, (SELECT COUNT(1) FROM adx.airdrop_submissions AS asub WHERE asub.referrer = '%s' AND asub.airdrop_id = %d) AS submissions FROM adx.codes WHERE wallet='%s' AND airdrop_id=%d LIMIT 1", db.Escape(req.Wallet), proto.AirdropId, db.Escape(req.Wallet), proto.AirdropId)
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	if len(rows) > 0 {
		promotion, err := getPromotionLink(proto, req.Wallet)
		if CheckErr(err, c) {
			log.Error(err.Error())
			return
		}
		code := token.Token(rows[0].Uint64(0))
		promotion.VerifyCode = code
        promotion.Submissions = rows[0].Uint64(1)

		c.JSON(http.StatusOK, promotion)
		return
	}
	var email = "NULL"
	if req.Email != "" && requireEmail > 0 {
		email = fmt.Sprintf("'%s'", db.Escape(req.Email))
	}
	_, _, err = db.Query("UPDATE adx.codes SET wallet='%s', email=%s, referrer='%s', `status`=1 WHERE id=%d AND `status`=0", db.Escape(req.Wallet), email, db.Escape(proto.Referrer), req.Code)
	//if strings.Contains(err.Error(), "Duplicate entry") {
	//	err = nil
	//}
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	_, _, err = db.Query(`INSERT INTO adx.airdrop_wallets (airdrop_id, wallet, referrer, email) VALUES (%d, '%s', '%s', %s) ON DUPLICATE KEY UPDATE email=VALUES(email)`, proto.AirdropId, db.Escape(req.Wallet), db.Escape(proto.Referrer), email)
	if err != nil {
		log.Error(err.Error())
	}
	promotion, err := getPromotionLink(proto, req.Wallet)
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, promotion)
}

func getPromotionLink(proto common.PromotionProto, wallet string) (promotion common.Promotion, err error) {
	promo := common.PromotionProto{
		Id:        proto.Id,
		UserId:    proto.UserId,
		AirdropId: proto.AirdropId,
		AdzoneId:  proto.AdzoneId,
		ChannelId: proto.ChannelId,
		Referrer:  wallet,
	}

	promoKey, err := common.EncodePromotion([]byte(Config.LinkSalt), promo)
	if err != nil {
		return promotion, err
	}
	link := fmt.Sprintf("%s/promo/%s", Config.BaseUrl, promoKey)
	shortURL, err := shorturl.Sina(link)
	if err == nil && shortURL != "" {
		link = shortURL
	}
	promotion = common.Promotion{
		Link: link,
		Key:  promoKey,
	}
	return promotion, nil
}
