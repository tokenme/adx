package media

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"github.com/tokenme/adx/common"
	"time"
)

func IndexMediaHander(c *gin.Context){
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsAdmin!=1,"is not admin",c){
		return
	}
	db := Service.Db
	id := c.Param("id")
	row,resut,err:=db.Query(`SELECT a.id,a.user_id,t.mobile,t.email,
	a.title,a.domain,a.intro,a.verified,a.online_status,a.verified_at,
	a.inserted_at,a.updated_at FROM adx.medias AS a 
	INNER JOIN adx.users AS t ON (a.user_id=t.id) 
	WHERE a.id=%s `,id)
	CheckErr(err,c)
	if len(row) == 0{
		c.JSON(http.StatusNotFound,"Not Find Media")
		return
	}
	Media := common.Media{
		Id:           row[0].Uint64(resut.Map("id")),
		UserId:       row[0].Uint64(resut.Map("user_id")),
		Title:        row[0].Str(resut.Map("title")),
		Domain:       row[0].Str(resut.Map("domain")),
		Intro:		  row[0].Str(resut.Map("intro")),
		Verified:     row[0].Uint(resut.Map("verified")),
		OnlineStatus: row[0].Uint(resut.Map("online_status")),
		Mobile: 	  row[0].Str(resut.Map("mobile")),
		Email: 		  row[0].Str(resut.Map("email")),
		Verified_at:  row[0].Time((resut.Map("verified_at")), time.Local),
		InsertedAt:   row[0].Time((resut.Map("inserted_at")), time.Local),
		UpdatedAt:    row[0].Time(resut.Map("updated_at"), time.Local),
	}

	c.JSON(http.StatusOK,Media)
}


func MediaInfoHandler(c *gin.Context){
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsAdmin!=1,"is not admin",c){
		return
	}
	db := Service.Db
	row,resut,err:=db.Query(`SELECT a.id,a.user_id,t.mobile,t.email,
	a.title,a.domain,a.intro,a.verified,a.online_status,a.verified_at,
	a.inserted_at,a.updated_at FROM adx.medias AS a 
	INNER JOIN adx.users AS t ON (a.user_id=t.id)`)
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
	c.JSON(http.StatusOK,info)
}
