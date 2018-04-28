package ad

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"time"
)

func ClickHandler(c *gin.Context) {
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
	err = AdClickQueue.NewClick(adKey)
	if CheckErr(err, c) {
		return
	}
	c.Redirect(http.StatusFound, ad.Url)
	return
}
