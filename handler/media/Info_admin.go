package media

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"github.com/tokenme/adx/common"
	"time"
)

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
	row,resut,err:=db.Query(`SELECT id,user_id,title,intro,online_status,inserted_at,updated_at FROM medias`)
	CheckErr(err,c)
	info:=[]common.Media{}
	for _,value:=range row{
		Media := common.Media{
			Id:           value.Uint64(resut.Map("id")),
			UserId:       value.Uint64(resut.Map("user_id")),
			Title:        value.Str(resut.Map("title")),
			OnlineStatus: value.Uint(resut.Map("online_status")),
			InsertedAt:   value.Time((resut.Map("inserted_at")), time.Local),
			UpdatedAt:    value.Time(resut.Map("updated_at"), time.Local),
		}
		info = append(info, Media)
		if err !=nil{
			break
		}
	}
	c.JSON(http.StatusOK,info)
}
