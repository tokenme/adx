package shorturl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SinaResponse struct {
	Urls []SinaResponseUrl `json:"urls"`
}

type SinaResponseUrl struct {
	Result bool   `json:"result"`
	Short  string `json:"url_short"`
	Long   string `json:"url_long"`
}

func Sina(link string) (short string, err error) {
	call := fmt.Sprintf("https://api.weibo.com/2/short_url/shorten.json?source=2849184197&url_long=%s", url.QueryEscape(link))
	resp, err := http.Get(call)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res SinaResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if len(res.Urls) == 0 || res.Urls[0].Short == "" {
		return "", errors.New("no response")
	}
	return res.Urls[0].Short, nil
}
