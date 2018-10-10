package promotion

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/adx/coins/eth"
	. "github.com/tokenme/adx/handler"
	"net/http"
)

type NewWalletResponse struct {
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

func NewWalletHandler(c *gin.Context) {
	privateKey, publicKey, err := eth.GenerateAccount()
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, NewWalletResponse{PublicKey: publicKey, PrivateKey: privateKey})
}
