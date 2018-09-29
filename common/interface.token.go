package common

import (
	"fmt"
	cmc "github.com/coincircle/go-coinmarketcap/types"
	"github.com/tokenme/adx/tools/ethplorer-api"
)

type Token struct {
	Address       string                `json:"address,omitempty"`
	Name          string                `json:"name,omitempty"`
	Symbol        string                `json:"symbol,omitempty"`
	Decimals      uint                  `json:"decimals,omitempty"`
	Protocol      string                `json:"protocol,omitempty"`
	Price         *ethplorer.TokenPrice `json:"price,omitempty"`
	Logo          uint                  `json:"logo,omitempty"`
	LogoAddress   string                `json:"logo_address,omitempty"`
	Summary       map[string]string     `json:"summary,omitempty"`
	Website       string                `json:"website,omitempty"`
	Blog          string                `json:"blog,omitempty"`
	Telegram      string                `json:"telegram,omitempty"`
	Facebook      string                `json:"facebook,omitempty"`
	Twitter       string                `json:"twitter,omitempty"`
	Whitepaper    string                `json:"whitepaper,omitempty"`
	Email         string                `json:"email,omitempty"`
	ClientIOS     string                `json:"client_ios,omitempty"`
	ClientAndroid string                `json:"client_android,omitempty"`
	Graph         cmc.TickerGraph       `json:"graph,omitempty"`
}

func (this Token) GetLogoAddress(cdn string) string {
	var addr string
	if this.Logo == 0 {
		addr = "default"
	} else if this.Address == "" {
		addr = "ethereum"
	} else {
		addr = this.Address
	}
	return fmt.Sprintf("%simg/token/%s.png", cdn, addr)
}

type TokenMarket struct {
	Id                string  `json:"id"`
	PriceUSD          float64 `json:"price_usd"`
	TotalSupply       float64 `json:"total_supply"`
	CirculatingSupply float64 `json:"available_supply"`
	MarketCapUSD      float64 `json:"market_cap_usd"`
	Volume24H         float64 `json:"24h_volume_usd"`
	PercentChange24H  float64 `json:"percent_change_24h"`
}
