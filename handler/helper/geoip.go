package helper

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"net"
	"net/http"
)

func GeoIPHandler(c *gin.Context) {
	clientIP := net.ParseIP(ClientIP(c))
	country, err := Service.GeoIP.Country(clientIP)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, country)
}
