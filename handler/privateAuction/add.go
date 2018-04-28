package privateAuction

import (
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"math/big"
	"net/http"
	"strings"
)

type AddRequest struct {
	AdzoneId  uint64               `form:"adzone_id" json:"adzone_id" bindinng:"required"`
	MediaId   uint64               `form:"media_id" json:"media_id" binding:"required"`
	Price     float64              `form:"price" json:"price" binding:"required"`
	Title     string               `form:"title" json:"title" binding:"required"`
	StartTime string               `form:"start_time" json:"start_time" binding:"required"`
	EndTime   string               `form:"end_time" json:"end_time" binding:"required"`
	Creatives []AddCreativeRequest `form:"creatives" json:"creatives" binding:"required"`
}

type AddCreativeRequest struct {
	Url string `form:"url" json:"url" binding:"required"`
	Img string `form:"img" json:"img" binding:"required"`
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

	if Check(user.IsAdvertiser != 1, "unauthorized", c) {
		return
	}
	balance, err := user.Balance(c, Service, Config)
	if CheckErr(err, c) {
		return
	}
	price := new(big.Int).SetUint64(uint64(req.Price * params.Ether))
	if balance.Cmp(price) == -1 {
		c.JSON(http.StatusOK, APIError{Code: NO_ENOUGH_BALANCE_ERROR, Msg: "no enough balance"})
		return
	}
	db := Service.Db
	_, ret, err := db.Query(`INSERT INTO adx.private_auctions (user_id, adzone_id, media_id, title, price, start_on, end_on) VALUES (%d, %d, %d, '%s', %.18f, '%s', '%s')`, user.Id, req.AdzoneId, req.MediaId, db.Escape(req.Title), req.Price, db.Escape(req.StartTime), db.Escape(req.EndTime))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	auctionId := ret.InsertId()
	var val []string
	for _, c := range req.Creatives {
		val = append(val, fmt.Sprintf("(%d, %d, %d, %d, '%s', '%s')", user.Id, auctionId, req.AdzoneId, req.MediaId, db.Escape(c.Url), db.Escape(c.Img)))
	}
	_, _, err = db.Query(`INSERT INTO adx.private_auction_creatives (user_id, auction_id, adzone_id, media_id, landing_page, img) VALUES %s`, strings.Join(val, ","))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
