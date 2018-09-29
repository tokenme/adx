package token

import (
	"github.com/gin-gonic/gin"
	cmc "github.com/miguelmota/go-coinmarketcap"
	. "github.com/tokenme/adx/handler"
	"net/http"
	"strings"
	"time"
)

func GraphHandler(c *gin.Context) {
	q := c.Query("q")
	var coinId string
	if q == "ETH" {
		coinId = "ETH"
	} else {
		q = strings.ToLower(q)
		db := Service.Db
		rows, _, err := db.Query(`SELECT symbol FROM adx.tokens WHERE address='%s' LIMIT 1`, db.Escape(q))
		if CheckErr(err, c) {
			return
		}
		if len(rows) == 0 {
			c.JSON(http.StatusOK, APIError{Code: NOTFOUND_ERROR, Msg: "not found"})
			return
		}

		row := rows[0]
		coinId = row.Str(0)
		coinId = strings.ToLower(coinId)
		coinId = strings.Replace(coinId, " ", "-", 0)
	}

	var (
		end   = time.Now()
		start = time.Now().AddDate(0, 0, -1)
	)
	if c.Query("start") != "" {
		start, _ = time.Parse("2006-01-02", c.Query("start"))
	}
	if c.Query("end") != "" {
		end, _ = time.Parse("2006-01-02", c.Query("end"))
	}

	options := &cmc.TickerGraphOptions{
		Symbol: coinId,
		Start:  start.Unix(),
		End:    end.Unix(),
	}
	graph, err := cmc.TickerGraph(options)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, graph)
}
