package token

import (
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"strings"
)

func GetHandler(c *gin.Context) {
	q := c.Query("q")
	if Check(q == "", "please use ERC20 token Name/Symbol/Address", c) {
		return
	}
	if q == "ETH" {
		token := common.Token{
			Name:        "Ethereum",
			Symbol:      "ETH",
			Decimals:    18,
			Protocol:    "ERC20",
			Logo:        1,
			Website:     "https://www.ethereum.org/",
			Blog:        "https://blog.ethereum.org/",
			Facebook:    "ethereumproject",
			Twitter:     "ethereumproject",
			Email:       "info@ethereum.org",
			LogoAddress: "https://www.ethereum.org/images/logos/ETHEREUM-ICON_Black.png",
		}
		c.JSON(http.StatusOK, token)
		return
	}
	db := Service.Db
	var (
		where        string
		addressQuery bool
	)
	if strings.HasPrefix(q, "0x") {
		addressQuery = true
		q = strings.ToLower(q)
		where = fmt.Sprintf("address='%s'", db.Escape(q))
	} else {
		escapedQ := db.Escape(q)
		where = fmt.Sprintf("name='%s' OR symbol='%s'", escapedQ, escapedQ)
	}
	rows, _, err := db.Query(`SELECT address, name, symbol, decimals, protocol, logo, summary, website, blog, telegram, facebook, twitter, whitepaper, email FROM tokenme.tokens WHERE %s LIMIT 1`, where)
	if CheckErr(err, c) {
		return
	}
	if len(rows) > 0 {
		row := rows[0]
		summary := make(map[string]string)
		json.Unmarshal([]byte(row.Str(6)), &summary)
		token := common.Token{
			Address:    row.Str(0),
			Name:       row.Str(1),
			Symbol:     row.Str(2),
			Decimals:   row.Uint(3),
			Protocol:   row.Str(4),
			Logo:       row.Uint(5),
			Summary:    summary,
			Website:    row.Str(7),
			Blog:       row.Str(8),
			Telegram:   row.Str(9),
			Facebook:   row.Str(10),
			Twitter:    row.Str(11),
			Whitepaper: row.Str(12),
			Email:      row.Str(13),
		}
		token.LogoAddress = token.GetLogoAddress(Config.CDNUrl)
		c.JSON(http.StatusOK, token)
		return
	}
	if !addressQuery {
		c.JSON(http.StatusOK, APIError{Code: NOTFOUND_ERROR, Msg: "not found"})
		return
	}
	geth := Service.Geth
	if geth == nil {
		c.JSON(http.StatusOK, APIError{Code: NOTFOUND_ERROR, Msg: "not found"})
		return
	}
	tokenCaller, err := eth.NewTokenCaller(ethcommon.HexToAddress(q), geth)
	if CheckErr(err, c) {
		return
	}
	tokenSymbol, err := tokenCaller.Symbol(nil)
	if CheckErr(err, c) {
		return
	}
	tokenDecimals, err := tokenCaller.Decimals(nil)
	if CheckErr(err, c) {
		return
	}
	tokenName, err := tokenCaller.Name(nil)
	if CheckErr(err, c) {
		return
	}
	token := common.Token{
		Address:  q,
		Name:     tokenName,
		Symbol:   tokenSymbol,
		Decimals: uint(tokenDecimals),
		Protocol: "ERC20",
	}
	_, _, err = db.Query(`INSERT IGNORE INTO tokenme.tokens (address, name, symbol, decimals, protocol) VALUES ('%s', '%s', '%s', %d, '%s')`, db.Escape(token.Address), db.Escape(token.Name), db.Escape(token.Symbol), token.Decimals, token.Protocol)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, token)
}
