package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	adServer "github.com/tokenme/adx/tools/ad/server"
	"github.com/tokenme/adx/tools/sqs"
	"github.com/tokenme/adx/tools/tracker"
	"net"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var (
	Service      *common.Service
	Config       common.Config
	EmailQueue   *sqs.EmailQueue
	AdClickQueue *sqs.AdClickQueue
	AdImpQueue   *sqs.AdImpQueue
	AdServer     *adServer.Server
	Tracker      *tracker.Tracker
)

func InitHandler(s *common.Service, c common.Config, ad *adServer.Server, emailQueue *sqs.EmailQueue, adClickQueue *sqs.AdClickQueue, adImpQueue *sqs.AdImpQueue, t *tracker.Tracker) {
	Service = s
	Config = c
	Tracker = t
	AdServer = ad
	EmailQueue = emailQueue
	AdClickQueue = adClickQueue
	AdImpQueue = adImpQueue
	raven.SetDSN(Config.SentryDSN)
}

type APIResponse struct {
	Msg string `json:"message,omitempty"`
}

type ErrorCode uint

const (
	BADREQUEST_ERROR        ErrorCode = 400
	INTERNAL_ERROR          ErrorCode = 500
	NOTFOUND_ERROR          ErrorCode = 404
	UNAUTHORIZED_ERROR      ErrorCode = 401
	INVALID_PASSWD_ERROR    ErrorCode = 409
	DUPLICATE_USER_ERROR    ErrorCode = 202
	UNACTIVATED_USER_ERROR  ErrorCode = 502
	DUPLICATE_MEDIA_ERROR   ErrorCode = 402
	UNVERIFIED_MEDIA_ERROR  ErrorCode = 405
	NO_ENOUGH_BALANCE_ERROR ErrorCode = 600
)

type APIError struct {
	Code ErrorCode `json:"code,omitempty"`
	Msg  string    `json:"message,omitempty"`
}

func (this APIError) Error() string {
	return fmt.Sprintf("CODE:%d, MSG:%s", this.Code, this.Msg)
}

func Check(flag bool, err string, c *gin.Context) (ret bool) {
	ret = flag
	if ret {
		_, file, line, _ := runtime.Caller(1)
		log.Error("[%s:%d]: %s", path.Base(file), line, err)
		c.JSON(http.StatusOK, APIError{Code: BADREQUEST_ERROR, Msg: err})
	}
	return
}
func CheckErr(err error, c *gin.Context) (ret bool) {
	ret = err != nil
	if ret {
		_, file, line, _ := runtime.Caller(1)
		log.Error("[%s:%d]: %s", path.Base(file), line, err)
		c.JSON(http.StatusOK, APIError{Code: BADREQUEST_ERROR, Msg: err.Error()})
	}
	return
}

func Uint64Value(val string, defaultVal uint64) (uint64, error) {
	if val == "" {
		return defaultVal, nil
	}

	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func Uint64NonZero(val string, err string) (uint64, error) {
	if val == "" {
		return 0, errors.New(err)
	}

	i, e := strconv.ParseUint(val, 10, 64)
	if e != nil {
		return 0, e
	}

	return i, nil
}

func ClientIP(c *gin.Context) string {
	if values, _ := c.Request.Header["X-Forwarded-For"]; len(values) > 0 {
		clientIP := values[0]
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}
	if values, _ := c.Request.Header["X-Real-Ip"]; len(values) > 0 {
		clientIP := strings.TrimSpace(values[0])
		if len(clientIP) > 0 {
			return clientIP
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func IP2Long(IpStr string) (int64, error) {
	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}

	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24-8*i)
	}

	return sum, nil
}

func Json(obj interface{}) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(obj)

	return buf.String()
}
