package ethplorer

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (this *Client) GetTokenInfo(tokenAddress string) (token Token, err error) {
	uri := fmt.Sprintf("/getTokenInfo/%s", tokenAddress)
	data, err := this.Exec(uri, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return
	}
	mp := make(map[string]interface{})
	err = json.Unmarshal(data, &mp)
	if err != nil {
		return
	}
	token.Price = parseTokenPrice(mp["price"])
	return token, nil
}

func parseTokenPrice(p interface{}) (t TokenPrice) {
	switch p.(type) {
	case bool:
		return
	case map[string]interface{}:
		m := p.(map[string]interface{})
		if v, found := m["rate"]; found {
			t.Rate, _ = strconv.ParseFloat(v.(string), 10)
		}
		if v, found := m["currency"]; found {
			t.Currency = v.(string)
		}
		if v, found := m["diff"]; found {
			t.Diff = v.(float64)
		}
		if v, found := m["diff7d"]; found {
			t.Diff7d = v.(float64)
		}
		if v, found := m["ts"]; found {
			t.Ts, _ = strconv.ParseInt(v.(string), 10, 64)
		}
		if v, found := m["marketCapUsd"]; found {
			t.MarketCapUsd, _ = strconv.ParseFloat(v.(string), 10)
		}
		if v, found := m["availableSupply"]; found {
			t.AvailableSupply, _ = strconv.ParseFloat(v.(string), 10)
		}
		if v, found := m["volume24h"]; found {
			t.Volume24h, _ = strconv.ParseFloat(v.(string), 10)
		}
	}
	return
}
