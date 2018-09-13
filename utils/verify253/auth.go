package verify253

import (
	. "github.com/tokenme/adx/handler"

	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/tools/zz253"
	"math/rand"
	"time"
	"net/http"
	"github.com/tokenme/adx/utils/twilio"
	"log"
)

func AuthSend(telephone string, c *gin.Context)  {
	db := Service.Db
	row, _, err := db.Query(fmt.Sprintf(`SELECT id FROM adx.telephone_codes WHERE telephone = '%s' AND type = 5 AND created >= '%s' LIMIT 1`, telephone, time.Now().Add(-1*time.Minute).Format("2006-01-02 15:04:05")))
	if Check(row != nil || err != nil, "请稍后1分钟后重试", c) {
		return
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
		Account:Config.Zz253.Account,
		Password:Config.Zz253.Password,
	}
	log.Println(apiReq.BaseRequest)
	_, err = zz253.Send(&apiReq)
	if CheckErr(err, c) {
		return
	}

	_, _, err = db.Query(`INSERT INTO adx.telephone_codes(telephone, type, 
	code, created) VALUES ('%s', %d, '%s', NOW()) ON
	DUPLICATE KEY UPDATE code=VALUES(code), created=VALUES(created)`, db.Escape(telephone), 5, db.Escape(code))
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}

func AuthVerification(Mobile string,Code string,Country uint) (twilio.AuthVerificationResponse,error) {
	db := Service.Db
	Res:=twilio.AuthVerificationResponse{Success:true,Message:"OK"}
	row, _, err := db.Query(`SELECT * FROM adx.telephone_codes WHERE telephone = '%s' AND code = '%s'`, Mobile, Code)
	if err !=nil || row == nil{
		Res.Success = false
		Res.Message = "验证码错误,请重新输入"
		return  Res,err
	}

	return Res,nil


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
