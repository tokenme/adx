package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

func ResetPwdVerifyHandler(c *gin.Context) {
	email := c.Query("email")
	if Check(email == "", "missing email", c) {
		return
	}

	code := c.Query("code")
	if Check(code == "", "missing verification code", c) {
		return
	}

	db := Service.Db
	rows, _, err := db.Query(`SELECT 1 FROM adx.users WHERE email='%s' AND reset_pwd_code='%s' AND reset_pwd_time >= DATE_SUB(NOW(), INTERVAL 2 HOUR)`, db.Escape(email), db.Escape(code))
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"Title":   "Error",
			"Message": "Wrong email or verification code expired",
		})
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("https://%s/#/reset-passwd/%s", c.Request.Host, code))
	return
}
