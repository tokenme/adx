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

	row,resut,err:=db.Query(`SELECT id,user_id,title,intro,online_status,inserted_at,updated_at FROM medias WHERE id=%s`,id)
	CheckErr(err,c)
	if len(row) == 0{
		c.JSON(http.StatusNotFound,"Not Find Media")
		return
	}
	Media := common.Media{
		Id:           row[0].Uint64(resut.Map("id")),
		UserId:       row[0].Uint64(resut.Map("user_id")),
		Title:        row[0].Str(resut.Map("title")),
		OnlineStatus: row[0].Uint(resut.Map("online_status")),
		InsertedAt:   row[0].Time((resut.Map("inserted_at")), time.Local),
		UpdatedAt:    row[0].Time(resut.Map("updated_at"), time.Local),
	}

	c.JSON(http.StatusOK,Media)
}
