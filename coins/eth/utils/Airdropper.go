package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tokenme/adx/coins/eth"
	"math/big"
)

func DeployAirdropper(auth *bind.TransactOpts, backend bind.ContractBackend, tokenAddress string) (common.Address, *types.Transaction, *eth.Airdropper, error) {
	return eth.DeployAirdropper(auth, backend, common.HexToAddress(tokenAddress))
}

func NewAirdropper(addr string, backend bind.ContractBackend) (*eth.Airdropper, error) {
	return eth.NewAirdropper(common.HexToAddress(addr), backend)
}

func AirdropperDrop(token *eth.Airdropper, opts *bind.TransactOpts, recipients []common.Address, tokenAmounts []*big.Int) (*types.Transaction, error) {
	return token.Airdrop(opts, recipients, tokenAmounts)
}
