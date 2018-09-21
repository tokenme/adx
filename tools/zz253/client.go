package zz253

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Request struct {
	MethodUrl string
	Params    map[string]interface{}

	IsDebug bool
}

func NewRequest(methodUrl string) *Request {
	return &Request{
		MethodUrl: methodUrl,
		Params:    make(map[string]interface{}),
	}
}

type Client struct {
	Account  string
	Password string

	IsTest  bool
	IsDebug bool
}

//create new client
func NewClient(account, password string, isTest bool) (c *Client) {
	c = &Client{
		Account:  account,
		Password: password,

		IsTest: isTest,
	}
	return
}

func (c *Client) Execute(req *Request) (string, error) {

	sysParams := make(map[string]interface{})
	sysParams["account"] = c.Account
	sysParams["password"] = c.Password
	for k, v := range req.Params {
		sysParams[k] = v
	}

	gatewayUrl := GATEWAY_URL
	if c.IsTest {
		gatewayUrl = TEST_GATEWAY_URL
	}
	if c.IsDebug {
		logger.Infof("request params: %s", Json(sysParams))
	}
	reqUrl := gatewayUrl + req.MethodUrl
	if c.IsDebug {
		logger.Infof("request url: %s", reqUrl)
	}
	response, err := http.DefaultClient.Post(reqUrl, "application/json; charset=UTF-8", strings.NewReader(Json(sysParams)))
	if err != nil {
		err = &Error{Code: "-1", Msg: err.Error()}
		logger.Error(err)
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = &Error{Code: "-1", Msg: err.Error()}
		logger.Error(err)
		return "", err
	}
	if c.IsDebug {
		logger.Infof("response: %s", string(body))
	}

	return string(body), nil
}
