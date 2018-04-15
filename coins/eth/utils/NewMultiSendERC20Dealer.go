package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tokenme/adx/coins/eth"
	"math/big"
)

func DeployMultiSendERC20Dealer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *eth.MultiSendERC20Dealer, error) {
	return eth.DeployMultiSendERC20Dealer(auth, backend)
}

func NewMultiSendERC20Dealer(addr string, backend bind.ContractBackend) (*eth.MultiSendERC20Dealer, error) {
	return eth.NewMultiSendERC20Dealer(common.HexToAddress(addr), backend)
}

func MultiSendERC20DealerTransfer(token *eth.MultiSendERC20Dealer, opts *bind.TransactOpts, tokenAddress string, dropper string, dealer string, commissionFee *big.Int, recipients []common.Address, tokenAmounts []*big.Int) (*types.Transaction, error) {
	return token.MultiSend(opts, common.HexToAddress(tokenAddress), common.HexToAddress(dropper), common.HexToAddress(dealer), commissionFee, recipients, tokenAmounts)
}
