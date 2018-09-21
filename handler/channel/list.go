package channel

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func ListGetHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT id, user_id, name FROM adx.channels WHERE user_id IN (0, %d) ORDER BY id ASC`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var channels []common.Channel
	for _, row := range rows {
		channels = append(channels, common.Channel{
			Id:     row.Uint64(0),
			UserId: row.Uint64(1),
			Name:   row.Str(2),
		})
	}
	c.JSON(http.StatusOK, channels)
}
