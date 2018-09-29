package token

import (
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	cmc "github.com/miguelmota/go-coinmarketcap"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"math"
	"math/big"
	"net/http"
	"strings"
)

func MarketHandler(c *gin.Context) {
	q := c.Query("q")
	var (
		coinId   string
		price    float64
		decimals *big.Int
	)
	if q == "ETH" {
		coinId = "ETH"
	} else {
		q = strings.ToLower(q)
		db := Service.Db
		rows, _, err := db.Query(`SELECT symbol, price, decimals FROM adx.tokens WHERE address='%s' LIMIT 1`, db.Escape(q))
		if CheckErr(err, c) {
			return
		}
		if len(rows) == 0 {
			c.JSON(http.StatusOK, APIError{Code: NOTFOUND_ERROR, Msg: "not found"})
			return
		}

		row := rows[0]
		price = row.ForceFloat(1)
		if row.Int(2) > 0 {
			decimals = big.NewInt(int64(math.Pow10(row.Int(2))))
		} else {
			decimals = big.NewInt(0)
		}

		coinId = row.Str(0)
		coinId = strings.ToLower(coinId)
		coinId = strings.Replace(coinId, " ", "-", 0)
	}

	options := &cmc.TickerOptions{
		Symbol: coinId,
	}
	coinTicker, err := cmc.Ticker(options)
	if err != nil && price > 0 {
		coin := common.TokenMarket{
			Id:       coinId,
			PriceUSD: price,
		}
		tokenCaller, err := eth.NewStandardTokenCaller(ethcommon.HexToAddress(q), Service.Geth)
		if CheckErr(err, c) {
			return
		}
		totalSupply, err := tokenCaller.TotalSupply(nil)
		if CheckErr(err, c) {
			return
		}
		if decimals.Cmp(big.NewInt(0)) == 0 {
			coin.TotalSupply = float64(totalSupply.Uint64())
		} else {
			coin.TotalSupply = float64(new(big.Int).Div(totalSupply, decimals).Uint64())
		}

		coin.MarketCapUSD = coin.TotalSupply * coin.PriceUSD
		c.JSON(http.StatusOK, coin)
		return
	} else if err != nil {
		c.JSON(http.StatusOK, common.TokenMarket{})
		return
	}
	coin := common.TokenMarket{
		Id:                coinId,
		TotalSupply:       coinTicker.TotalSupply,
		CirculatingSupply: coinTicker.CirculatingSupply,
	}
	if quote, found := coinTicker.Quotes["USD"]; found {
		coin.PriceUSD = quote.Price
		coin.MarketCapUSD = quote.MarketCap
		coin.Volume24H = quote.Volume24H
		coin.PercentChange24H = quote.PercentChange24H
	}
	redisMasterConn := Service.Redis.Master.Get()
	defer redisMasterConn.Close()
	redisMasterConn.Do("SETEX", fmt.Sprintf("coinprice-%s", strings.ToLower(q)), 60*60, coin.PriceUSD)
	c.JSON(http.StatusOK, coin)
	return
}
