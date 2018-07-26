package common

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/tokenme/adx/utils"
	"github.com/tokenme/adx/utils/binary"
	"math"
	"math/big"
	"time"
)

type Airdrop struct {
	Id            uint64    `json:"id"`
	UserId        uint64    `json:"user_id,omitempty"`
	Title         string    `json:"title,omitempty"`
	Token         Token     `json:"token"`
	OnlineStatus  int       `json:"online_status"`
	Budget        uint64    `json:"budget"`
	CommissionFee uint64    `json:"commission_fee"`
	GiveOut       uint64    `json:"give_out"`
	Bonus         uint      `json:"bonus"`
	TelegramGroup string    `json:"telegram_group"`
	TelegramBot   string    `json:"telegram_bot"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	InsertedAt    time.Time `json:"inserted_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AirdropPromotion struct {
	AirdropId   uint64 `json:"airdrop_id"`
	ShareWallet string `json:"share_wallet"`
}

func (this Airdrop) TotalTokenNeeded() *big.Int {
	tokenDecimals := big.NewInt(int64(math.Pow10(int(this.Token.Decimals))))
	budget := new(big.Int).Mul(new(big.Int).SetUint64(this.Budget), tokenDecimals)
	submissions := new(big.Int).SetUint64(this.Budget / this.GiveOut)
	bonus := new(big.Int).Mul(submissions, big.NewInt(int64(this.Bonus)))
	return new(big.Int).Add(budget, bonus)
}

func (this Airdrop) TotalETHNeeded(config AirdropConfig) *big.Int {
	submissions := new(big.Int).SetUint64(this.Budget / this.GiveOut)
	commissionFee := new(big.Int).Mul(submissions, new(big.Int).SetUint64(this.CommissionFee))
	gasFee := new(big.Int).Mul(big.NewInt(config.GasPrice), big.NewInt(config.GasLimit))
	contractFee := new(big.Int).Mul(big.NewInt(config.DealerContractGasPrice), big.NewInt(config.DealerContractGasLimit))
	gas := new(big.Int).Add(gasFee, contractFee)
	gwei := new(big.Int).Add(gas, commissionFee)
	return new(big.Int).Mul(gwei, big.NewInt(params.Shannon))
}

func EncodeAirdropPromotion(key []byte, proto AirdropPromotion) (string, error) {
	enc := binary.NewEncoder()
	enc.Encode(proto)
	return utils.AESEncryptBytes(key, enc.Buffer())
}

func DecodeAirdropPromotion(key []byte, cryptoText string) (proto AirdropPromotion, err error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return proto, err
	}
	dec := binary.NewDecoder()
	dec.SetBuffer(data)
	dec.Decode(&proto)
	return proto, nil
}
