package adzone

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func SizeListHandler(c *gin.Context) {
	db := Service.Db
	rows, _, err := db.Query(`SELECT id, width, height FROM adx.sizes`)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	var sizes []common.Size
	for _, row := range rows {
		size := common.Size{
			Id:     row.Uint(0),
			Width:  row.Uint(1),
			Height: row.Uint(2),
		}
		sizes = append(sizes, size)
	}
	c.JSON(http.StatusOK, sizes)
}
