package media

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net"
	"net/http"
	"strings"
)

type VerifyRequest struct {
	Id uint64 `form:"id" json:"id" binding:"required"`
}

func VerifyHandler(c *gin.Context) {
	var req VerifyRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if Check(user.IsPublisher != 1, "unauthorized", c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT domain, salt FROM adx.medias WHERE id=%d AND user_id=%d`, req.Id, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(len(rows) == 0, "unauthorized", c) {
		return
	}
	media := common.Media{
		Id:       req.Id,
		Domain:   rows[0].Str(0),
		Identity: rows[0].Str(1),
	}
	media = media.Complete()

	var verified bool
	resp, err := http.Head(media.VerifyURL)
	if err != nil || resp.StatusCode != 200 {
		domainName, err := net.LookupCNAME(fmt.Sprintf("%s.", media.VerifyDNS))
		if err != nil {
			fmt.Println(err.Error())
		}
		if strings.TrimRight(domainName, ".") == media.DNSValue {
			verified = true
		}
	} else {
		verified = true
	}

	if verified {
		_, _, err = db.Query(`UPDATE adx.medias SET verified=1, verified_at=NOW() WHERE id=%d`, media.Id)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}
	c.JSON(http.StatusOK, APIError{Code: UNVERIFIED_MEDIA_ERROR, Msg: "unverified"})
	return

}
