package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tokenme/adx/coins/eth"
	"math/big"
)

func NewToken(addr string, backend bind.ContractBackend) (*eth.Token, error) {
	return eth.NewToken(common.HexToAddress(addr), backend)
}

func BalanceOfToken(geth *ethclient.Client, tokenAddress string, wallet string) (*big.Int, error) {
	token, err := eth.NewToken(common.HexToAddress(tokenAddress), geth)
	if err != nil {
		return nil, err
	}
	return TokenBalanceOf(token, wallet)
}

func TokenBalanceOf(token *eth.Token, wallet string) (*big.Int, error) {
	return token.BalanceOf(nil, common.HexToAddress(wallet))
}

func Transfer(token *eth.Token, opts *bind.TransactOpts, _to string, _value *big.Int) (*types.Transaction, error) {
	return token.Transfer(opts, common.HexToAddress(_to), _value)
}

func TokenDecimal(token *eth.Token, opts *bind.CallOpts) (int, error) {
	d, err := token.Decimals(opts)
	if err != nil {
		return 0, err
	}
	return int(d), nil
}
