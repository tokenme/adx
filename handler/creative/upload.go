package creative

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/tools/s3"
	"github.com/tokenme/adx/utils"
	"net/http"
)

func UploadHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)

	file, _, err := c.Request.FormFile("file")
	if CheckErr(err, c) {
		return
	}
	uuid, _ := utils.UUID()
	key := fmt.Sprintf("%d-%s", user.Id, uuid)

	buf := new(bytes.Buffer)
	var w = bufio.NewWriter(buf)
	w.ReadFrom(file)
	w.Flush()
	_, err = s3.Upload(Config.S3, Config.S3.AdBucket, Config.S3.CreativePath, key, buf.Bytes())
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": key, "url": fmt.Sprintf("%s/%s/%s", Config.CreativeCDN, Config.S3.CreativePath, key)})
}
