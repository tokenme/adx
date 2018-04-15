package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tokenme/adx/coins/eth"
	"math/big"
)

func NewStandardToken(addr string, backend bind.ContractBackend) (*eth.StandardToken, error) {
	return eth.NewStandardToken(common.HexToAddress(addr), backend)
}

func StandardTokenAllowance(token *eth.StandardToken, wallet string, contract string) (*big.Int, error) {
	return token.Allowance(nil, common.HexToAddress(wallet), common.HexToAddress(contract))
}

func StandardTokenBalanceOf(token *eth.StandardToken, wallet string) (*big.Int, error) {
	return token.BalanceOf(nil, common.HexToAddress(wallet))
}

func StandardTokenApprove(token *eth.StandardToken, opts *bind.TransactOpts, spender string, value *big.Int) (*types.Transaction, error) {
	return token.Approve(opts, common.HexToAddress(spender), value)
}
