package media

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"strconv"
	"time"
)

func IndexMediaHander(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsAdmin != 1, "is not admin", c) {
		return
	}
	db := Service.Db
	id := c.DefaultQuery("id", "1")
	row, resut, err := db.Query(`SELECT a.id,a.user_id,t.mobile,t.email,
	a.title,a.domain,a.intro,a.verified,a.online_status,a.verified_at,
	a.inserted_at,a.updated_at FROM adx.medias AS a 
	INNER JOIN adx.users AS t ON (a.user_id=t.id) 
	WHERE a.id=%s `, id)
	CheckErr(err, c)
	if len(row) == 0 {
		c.JSON(http.StatusNotFound, "Not Find Media")
		return
	}
	Media := common.Media{
		Id:           row[0].Uint64(resut.Map("id")),
		UserId:       row[0].Uint64(resut.Map("user_id")),
		Title:        row[0].Str(resut.Map("title")),
		Domain:       row[0].Str(resut.Map("domain")),
		Intro:        row[0].Str(resut.Map("intro")),
		Verified:     row[0].Uint(resut.Map("verified")),
		OnlineStatus: row[0].Uint(resut.Map("online_status")),
		Mobile:       row[0].Str(resut.Map("mobile")),
		Email:        row[0].Str(resut.Map("email")),
		Verified_at:  row[0].Time((resut.Map("verified_at")), time.Local),
		InsertedAt:   row[0].Time((resut.Map("inserted_at")), time.Local),
		UpdatedAt:    row[0].Time(resut.Map("updated_at"), time.Local),
	}

	c.JSON(http.StatusOK, Media)
}

const Limit = 15
const one = 1

func MediaInfoHandler(c *gin.Context){
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsAdmin!=1,"Is Not Admin",c){
		return
	}
	page,err:=strconv.Atoi(c.DefaultQuery("page","1"))
	CheckErr(err,c)
	if page<=0 {
		page = 1
	}
	index:=(page-one)*Limit
	db := Service.Db
	row,resut,err:=db.Query(`SELECT a.id,a.user_id,t.mobile,t.email,
	a.title,a.domain,a.intro,a.verified,a.online_status,a.verified_at,
	a.inserted_at,a.updated_at FROM adx.medias AS a 
	INNER JOIN adx.users AS t ON (a.user_id=t.id) LIMIT %d OFFSET %d`,Limit,index)
	CheckErr(err,c)
	info:=[]common.Media{}
	for _,value:=range row{
		Media := common.Media{
			Id:           value.Uint64(resut.Map("id")),
			UserId:       value.Uint64(resut.Map("user_id")),
			Title:        value.Str(resut.Map("title")),
			Domain:       value.Str(resut.Map("domain")),
			Intro:		  value.Str(resut.Map("intro")),
			Verified:     value.Uint(resut.Map("verified")),
			OnlineStatus: value.Uint(resut.Map("online_status")),
			Mobile: 	  value.Str(resut.Map("mobile")),
			Email: 		  value.Str(resut.Map("email")),
			Verified_at:  value.Time((resut.Map("verified_at")), time.Local),
			InsertedAt:   value.Time((resut.Map("inserted_at")), time.Local),
			UpdatedAt:    value.Time(resut.Map("updated_at"), time.Local),
		}
		info = append(info, Media)
	}
	row,_,err=db.Query(`SELECT COUNT(*) FROM adx.medias`)
	CheckErr(err,c)
	c.JSON(http.StatusOK,gin.H{
		"Total":row[0].Uint(0),
		"Data":info,
	})
}