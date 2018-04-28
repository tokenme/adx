package ad

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ShowResponse struct {
	Link string `json:"link"`
	Img  string `json:"img"`
	Imp  string `json:"imp"`
}

func ShowHandler(c *gin.Context) {
	encodedAdzoneId := c.Query("a")
	adzoneId, err := common.DecodeAdzoneId([]byte(Config.LinkSalt), encodedAdzoneId)
	if CheckErr(err, c) {
		return
	}
	ad, err := AdServer.PrivateAuction(adzoneId)
	if CheckErr(err, c) {
		return
	}
	uuid, err := c.Cookie("uuid")
	if err != nil {
		uuid, _ = utils.Salt()
		du := int(time.Hour * 24 * 30 / time.Second)
		c.SetCookie("uuid", uuid, du, "/", Config.CookieDomain, true, true)
	}
	bwType, _ := strconv.ParseUint(c.Query("bw_type"), 10, 64)
	ipInfo := ClientIP(c)
	ipNumber, _ := IP2Long(ipInfo)
	var (
		countryId   uint
		countryName string
	)
	clientIP := net.ParseIP(ipInfo)
	country, err := Service.GeoIP.Country(clientIP)
	if err == nil && country != nil {
		countryId = country.Country.GeoNameID
		if cname, found := country.Country.Names["en"]; found {
			countryName = cname
		}
	}
	httpReq := c.Request

	env := common.AdEnv{
		Cookie:         uuid,
		URL:            c.Query("url"),
		Referrer:       c.Query("referrer"),
		ScreenSize:     fmt.Sprintf("%sx%s", c.Query("sw"), c.Query("sh")),
		AdSize:         fmt.Sprintf("%sx%s", c.Query("adw"), c.Query("adh")),
		OsName:         c.Query("os_name"),
		OsVersion:      c.Query("os_ver"),
		BrowserName:    c.Query("bw_name"),
		BrowserVersion: c.Query("bw_ver"),
		BrowserType:    uint(bwType),
		UserAgent:      httpReq.UserAgent(),
		CountryId:      countryId,
		CountryName:    countryName,
		IP:             ipNumber,
	}
	ad.Env = env
	resp := ShowResponse{
		Link: ad.GetLink(Config),
		Img:  ad.GetImgUrl(Config),
		Imp:  ad.GetImpUrl(Config),
	}
	js, _ := json.Marshal(resp)
	contentType := "Content-Type: application/javascript; charset=utf-8"
	content := fmt.Sprintf("window.TMM(%s);", string(js))
	c.Data(http.StatusOK, contentType, []byte(content))
	return

}
