package ethplorer

import (
	"encoding/json"
	"fmt"
)

func (this *Client) GetAddressInfo(address string, tokenAddress string) (addressInfo Address, err error) {
	uri := fmt.Sprintf("/getAddressInfo/%s", address)
	var params map[string]string
	if tokenAddress != "" {
		params["token"] = tokenAddress
	}
	data, err := this.Exec(uri, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &addressInfo)
	if err != nil {
		return
	}
	if len(addressInfo.Tokens) > 0 {
		mp := make(map[string]interface{})
		err = json.Unmarshal(data, &mp)
		if err != nil {
			return
		}
		switch mp["tokens"].(type) {
		case []interface{}:
			for idx, v := range mp["tokens"].([]interface{}) {
				addressInfo.Tokens[idx].Token.Price = parseTokenPrice(v.(map[string]interface{})["price"])
			}
		}
	}
	return addressInfo, nil
}
