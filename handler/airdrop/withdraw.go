package airdrop

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/adx/coins/eth"
	ethutils "github.com/tokenme/adx/coins/eth/utils"
	"github.com/tokenme/adx/common"
	. "github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/utils"
	"math/big"
	"net/http"
)

type WithdrawRequest struct {
	AirdropId   uint64  `form:"airdrop_id" json:"airdrop_id" binding:"required"`
	Wallet      string  `form:"wallet" json:"wallet" binding:"required"`
	GasPrice    uint64  `form:"gas_price" json:"gas_price" binding:"required"`
	TokenAmount float64 `form:"token_amount" json:"token_amount"`
	Ether       float64 `form:"eth" json:"eth"`
	Passwd      string  `form:"passwd" json:"passwd" binding:"required"`
}

func WithdrawHandler(c *gin.Context) {
	var req WithdrawRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.TokenAmount > 0 && req.Ether > 0, "please transfer token and ether seperatly", c) {
		return
	}
	userContext, exists := c.Get("USER")
	if Check(!exists, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if Check(user.IsPublisher == 0 && user.IsAdmin == 0, "invalid permission", c) {
		return
	}
	db := Service.Db
	query := `SELECT
	uw.salt,
	uw.is_private,
	uw.passwd,
	u.passwd,
	u.salt
FROM adx.user_wallets AS uw
INNER JOIN tokenme.users AS u ON (u.id=uw.user_id)
WHERE uw.user_id=%d AND uw.is_main=1`
	rows, _, err := db.Query(query, user.Id)
	if CheckErr(err, c) {
		return
	}

	row := rows[0]
	salt := row.Str(0)
	isPrivate := row.Uint(1)
	password := row.Str(2)
	userPassword := row.Str(3)
	userSalt := row.Str(4)
	if Check(isPrivate != 1, "no privte key provided", c) {
		return
	}
	var validPassword bool
	if password == "" {
		passwdSha1 := utils.Sha1(fmt.Sprintf("%s%s%s", userSalt, req.Passwd, userSalt))
		validPassword = passwdSha1 == userPassword
	} else {
		passwdSha1 := utils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Passwd, salt))
		validPassword = passwdSha1 == password
	}
	if Check(!validPassword, "invalid password", c) {
		return
	}

	rows, _, err = db.Query(`SELECT t.address, t.name, t.symbol, t.decimals, t.protocol, a.wallet, a.salt FROM adx.airdrops AS a INNER JOIN adx.tokens AS t ON (t.address=a.token_address) WHERE a.id=%d AND a.wallet_val_t=0 AND t.protocol="ERC20" AND (a.user_id=%d OR EXISTS (SELECT 1 FROM adx.users WHERE id=%d AND is_admin=1 LIMIT 1)) LIMIT 1`, req.AirdropId, user.Id, user.Id)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(len(rows) == 0, "missing token", c) {
		return
	}
	token := common.Token{
		Address:  rows[0].Str(0),
		Name:     rows[0].Str(1),
		Symbol:   rows[0].Str(2),
		Decimals: rows[0].Uint(3),
		Protocol: rows[0].Str(4),
	}
	airdropWallet := rows[0].Str(5)
	airdropSalt := rows[0].Str(6)
	privateKey, err := utils.AddressDecrypt(airdropWallet, airdropSalt, Config.TokenSalt)
	if CheckErr(err, c) {
		return
	}
	publicKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if CheckErr(err, c) {
		return
	}
	var (
		totalTokens *big.Int
		totalEther  *big.Int
	)
	if token.Decimals >= 4 {
		totalTokensForSave := new(big.Int).SetUint64(uint64(req.TokenAmount * float64(utils.Pow40.Uint64())))
		totalTokens = new(big.Int).Mul(totalTokensForSave, utils.Pow10(int(token.Decimals)))
		totalTokens = new(big.Int).Div(totalTokens, utils.Pow40)
	} else {
		totalTokensForSave := new(big.Int).SetUint64(uint64(req.TokenAmount))
		totalTokens = new(big.Int).Mul(totalTokensForSave, utils.Pow10(int(token.Decimals)))
	}
	totalEtherForSave := new(big.Int).SetUint64(uint64(req.Ether * float64(utils.Pow40.Uint64())))
	totalEther = new(big.Int).Mul(totalEtherForSave, big.NewInt(params.Ether))
	totalEther = new(big.Int).Div(totalEther, utils.Pow40)

	if totalEther.Cmp(big.NewInt(0)) != 1 {
		totalEther = nil
	}
	transactor := eth.TransactorAccount(privateKey)
	nonce, err := eth.PendingNonce(Service.Geth, c, publicKey)
	if CheckErr(err, c) {
		return
	}
	gasPrice := new(big.Int).Mul(new(big.Int).SetUint64(req.GasPrice), big.NewInt(params.Shannon))
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 210000,
		Value:    totalEther,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	var tx *types.Transaction
	if totalTokens.Cmp(big.NewInt(0)) == 1 {
		tokenHandler, err := ethutils.NewToken(token.Address, Service.Geth)
		if CheckErr(err, c) {
			return
		}
		tx, err = ethutils.Transfer(tokenHandler, transactor, req.Wallet, totalTokens)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
	} else {
		tx, err = eth.Transfer(transactor, Service.Geth, c, req.Wallet)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
	}

	c.JSON(http.StatusOK, APIResponse{Msg: tx.Hash().Hex()})
}
