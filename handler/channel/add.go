package channel

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type AddRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
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
	_, _, err := db.Query(`INSERT IGNORE INTO adx.channels (user_id, name) VALUES (%d, '%s')`, user.Id, db.Escape(req.Name))
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
