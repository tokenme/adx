package auth

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils/twilio"
	"github.com/tokenme/adx/tools/zz253"
	"net/http"
	"strings"
	"fmt"
	"time"
	"math/rand"
	"bytes"
)

type SendRequest struct {
	Mobile  string `form:"mobile" json:"mobile" binding:"required"`
	Country uint   `form:"country" json:"country" binding:"required"`
}

func SendHandler(c *gin.Context) {
	var req SendRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	//fmt.Printf("req: %s\n", Json(req))
	mobile := strings.Replace(req.Mobile, " ", "", 0)
	if req.Country == 86 {
		err,msg := ChinaSend(mobile)
		if Check(err!=nil||msg!="",msg,c){
			return
		}
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}
	ret, err := twilio.AuthSend(Config.TwilioToken, mobile, req.Country)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(!ret.Success, ret.Message, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}


func ChinaSend(telephone string) (error, string){
	db := Service.Db
	Config := Config
	row, _, err := db.Query(fmt.Sprintf(`SELECT id FROM adx.telephone_codes WHERE telephone = '%s' AND type = 5 AND created >= '%s' LIMIT 1`, telephone, time.Now().Add(-1*time.Minute).Format("2006-01-02 15:04:05")))
	if row != nil || err != nil {
		return err, "请稍后1分钟后重试"
	}
	code := "0123456789"
	code = Shuffle(code)
	code = Shuffle(code)
	code = code[0:4]
	apiReq := zz253.SendRequest{}
	apiReq.Phones = []string{telephone}
	Content := fmt.Sprintf("【TOKENMAMA】您的注册验证码是: %s", code)
	apiReq.Content = Content
	apiReq.BaseRequest = zz253.BaseRequest{
		Account:  Config.Zz253.Account,
		Password: Config.Zz253.Password,
	}
	_, err = zz253.Send(&apiReq)
	if err != nil {
		return err, "短信发送失败"
	}
	_, _, err = db.Query(`INSERT INTO adx.telephone_codes(telephone, type, 
	code, created) VALUES ('%s', %d, '%s', NOW()) ON
	DUPLICATE KEY UPDATE code=VALUES(code), created=VALUES(created)`, db.Escape(telephone), 5, db.Escape(code))
	if err != nil {
		return err, "数据插入失败"
	}
	return nil, ""

}

func Shuffle(s string) (ret string) {
	var buffer bytes.Buffer
	r := rand.New(rand.NewSource(time.Now().Unix()))
	runes := []rune(s)
	perm := r.Perm(len(runes))
	for _, randIndex := range perm {
		buffer.WriteString(string(runes[randIndex]))
	}
	return buffer.String()

}