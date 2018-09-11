package media

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"github.com/ziutek/mymysql/mysql"
	"net/http"
	"net/url"
	"time"
)

type AddRequest struct {
	Title  string `form:"title" json:"title" binding:"required"`
	Domain string `form:"domain" json:"domain" binding:"required"`
	ImgUrl    string `form:"placeholder_img" json:"placeholder_img" binding:"required"`

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

	if Check(user.IsPublisher != 1 && user.IsAdvertiser != 1 && user.IsAdmin != 1, "unauthorized", c) {
		return
	}
	title := utils.Normalize(req.Title)
	domain := utils.Normalize(req.Domain)
	ImgUrl := utils.Normalize(req.ImgUrl)
	parsedUrl, err := url.Parse(domain)
	if CheckErr(err, c) {
		return
	}
	domain = fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
	identity, _ := utils.Salt()
	db := Service.Db
	_, ret, err := db.Query(`INSERT INTO adx.medias (user_id, title, domain, url, salt) VALUES (%d, '%s', '%s', '%s', '%s')`, user.Id, db.Escape(title), db.Escape(domain), db.Escape(ImgUrl), db.Escape(identity))
	if err != nil && err.(*mysql.Error).Code == mysql.ER_DUP_ENTRY {
		c.JSON(http.StatusOK, APIError{Code: DUPLICATE_MEDIA_ERROR, Msg: "media title or domain already exists"})
		return
	}
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	mediaId := ret.InsertId()
	media := common.Media{
		Id:         mediaId,
		Title:      title,
		Domain:     domain,
		ImgUrl:      ImgUrl,
		Identity:   identity,
		InsertedAt: time.Now(),
		UpdatedAt:  time.Now(),
	}
	media = media.Complete()
	c.JSON(http.StatusOK, media)
}
