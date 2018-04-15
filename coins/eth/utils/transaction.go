package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TransactionByHash(client *ethclient.Client, ctx context.Context, txHashHex string) (tx *types.Transaction, isPending bool, err error) {
	return client.TransactionByHash(ctx, common.HexToHash(txHashHex))
}

func TransactionReceipt(client *ethclient.Client, ctx context.Context, txHashHex string) (*types.Receipt, error) {
	return client.TransactionReceipt(ctx, common.HexToHash(txHashHex))
}
