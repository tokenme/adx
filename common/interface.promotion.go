package common

import (
	"github.com/tokenme/adx/utils"
	"github.com/tokenme/adx/utils/binary"
	"github.com/tokenme/adx/utils/token"
	"time"
)

type AirdropSubmission struct {
	Id             uint64
	Airdrop        *Airdrop
	Proto          PromotionProto
	Wallet         string
	PromoterWallet string
	Tx             string
}

type Promotion struct {
	Id              uint64      `json:"id"`
	UserId          uint64      `json:"user_id"`
	Airdrop         *Airdrop    `json:"airdrop"`
	AdzoneId        uint64      `json:"adzone_id"`
	ChannelId       uint64      `json:"channel_id"`
	AdzoneName      string      `json:"adzone_name"`
	ChannelName     string      `json:"channel_name"`
	Link            string      `json:"link"`
	Key             string      `json:"key"`
	SelfSubmissions uint64      `json:"self_submissions"`
	Submissions     uint64      `json:"submissions"`
	VerifyCode      token.Token `json:"verify_code,omitempty"`
	Inserted        time.Time   `json:"inserted"`
}

type PromotionStats struct {
	Pv            uint64    `json:"pv"`
	Submissions   uint64    `json:"submissions"`
	Transactions  uint64    `json:"transactions"`
	GiveOut       uint64    `json:"give_out"`
	Bonus         uint64    `json:"bonus"`
	CommissionFee uint64    `json:"commission_fee"`
	Decimals      uint      `json:"decimals"`
	RecordOn      time.Time `json:"record_on"`
}

type PromotionStatsWithSummary struct {
	Summary PromotionStats   `json:"summary"`
	Stats   []PromotionStats `json:"stats"`
}

type PromotionProto struct {
	Id        uint64 `json:"id"`
	UserId    uint64 `json:"user_id"`
	AirdropId uint64 `json:"airdrop_id"`
	AdzoneId  uint64 `json:"adzone_id"`
	ChannelId uint64 `json:"channel_id"`
	Referrer  string `json:"referrer"`
}

func EncodePromotion(key []byte, proto PromotionProto) (string, error) {
	enc := binary.NewEncoder()
	enc.Encode(proto)
	return utils.AESEncryptBytes(key, enc.Buffer())
}

func DecodePromotion(key []byte, cryptoText string) (proto PromotionProto, err error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return proto, err
	}
	dec := binary.NewDecoder()
	dec.SetBuffer(data)
	dec.Decode(&proto)
	return proto, nil
}
