package ethplorer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const GATEWAY = "https://api.ethplorer.io"

type Client struct {
	key string
}

func NewClient(key string) *Client {
	return &Client{key: key}
}

func (this *Client) Exec(uri string, params map[string]string) ([]byte, error) {
	values := url.Values{}
	values.Set("apiKey", this.key)
	if params != nil {
		for k, v := range params {
			values.Set(k, v)
		}
	}
	resp, err := http.Get(fmt.Sprintf("%s%s?%s", GATEWAY, uri, values.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
