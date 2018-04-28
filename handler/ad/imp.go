package ad

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"time"
)

func ImpHandler(c *gin.Context) {
	key := c.Param("key")
	ad, err := common.DecodeAd([]byte(Config.LinkSalt), key)
	if CheckErr(err, c) {
		return
	}
	ad.LogTime = time.Now().Unix()
	adKey, err := common.EncodeAd([]byte(Config.LinkSalt), ad)
	if CheckErr(err, c) {
		return
	}
	err = AdImpQueue.NewImp(adKey)
	if CheckErr(err, c) {
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("%simg/_.gif", Config.CDNUrl))
	return
}
