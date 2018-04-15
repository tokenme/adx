package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func ActivateHandler(c *gin.Context) {
	email := c.Query("email")
	if Check(email == "", "missing email", c) {
		return
	}

	activationCode := c.Query("activation_code")
	if Check(activationCode == "", "missing activation_code", c) {
		return
	}

	db := Service.Db
	_, ret, err := db.Query(`UPDATE adx.users SET active = 1 WHERE email='%s' AND activation_code='%s' AND activate_code_time >= DATE_SUB(NOW(), INTERVAL 2 HOUR)`, db.Escape(email), db.Escape(activationCode))
	if CheckErr(err, c) {
		return
	}
	if ret.AffectedRows() == 0 {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"Title":   "Error",
			"Message": "Wrong email or activation code expired",
		})
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("https://%s", c.Request.Host))
	return
}
