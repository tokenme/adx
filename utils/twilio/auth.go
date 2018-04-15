package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	AuthSendGatway         = "https://api.authy.com/protected/json/phones/verification/start"
	AuthVerficationGateway = "https://api.authy.com/protected/json/phones/verification/check"
)

type AuthSendResponse struct {
	Carrier     string `json:"carrier"`
	IsCellphone bool   `json:"is_cellphone"`
	Message     string `json:"message"`
	Expires     int64  `json:"seconds_to_expire"`
	Uuid        string `json:"uuid"`
	Success     bool   `json:"success"`
}

type AuthVerificationResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func AuthSend(key string, mobile string, country uint) (ret AuthSendResponse, err error) {
	value := url.Values{}
	value.Add("api_key", key)
	value.Add("phone_number", mobile)
	value.Add("country_code", strconv.FormatUint(uint64(country), 10))
	value.Add("via", "sms")
	resp, err := http.PostForm(AuthSendGatway, value)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(body, &ret)
	return
}

func AuthVerification(key string, mobile string, country uint, code string) (ret AuthVerificationResponse, err error) {
	value := url.Values{}
	value.Add("api_key", key)
	value.Add("phone_number", mobile)
	value.Add("country_code", strconv.FormatUint(uint64(country), 10))
	value.Add("verification_code", code)
	uri := fmt.Sprintf("%s?%s", AuthVerficationGateway, value.Encode())
	resp, err := http.DefaultClient.Get(uri)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(body, &ret)
	return
}
