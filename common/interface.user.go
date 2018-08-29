package common

import (
	"context"
	"errors"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	//"github.com/mkideal/log"
	"github.com/tokenme/adx/coins/eth"
	"github.com/tokenme/adx/utils"
	"math/big"
)

type User struct {
	Id             uint64        `json:"id,omitempty"`
	Mobile         string        `json:"mobile,omitempty"`
	CountryCode    uint          `json:"country_code,omitempty"`
	Email          string        `json:"email,omitempty"`
	ActivationCode string        `json:"-"`
	ResetPwdCode   string        `json:"-"`
	Telegram       *TelegramUser `json:"telegram,omitempty"`
	ShowName       string        `json:"showname,omitempty"`
	Avatar         string        `json:"avatar,omitempty"`
	Salt           string        `json:"-"`
	Password       string        `json:"-"`
	Wallet         string        `json:"wallet,omitempty"`
	IsAdmin        uint          `json:"is_admin"`
	IsPublisher    uint          `json:"is_publisher,omitempty"`
	IsAdvertiser   uint          `json:"is_advertiser,omitempty"`
}

type TelegramUser struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Avatar    string `json:"avatar"`
}

func (this User) GetShowName() string {
	if this.Telegram != nil && (this.Telegram.Firstname != "" && this.Telegram.Lastname != "" || this.Telegram.Username != "") {
		if this.Telegram.Username != "" {
			return this.Telegram.Username
		}
		return fmt.Sprintf("%s %s", this.Telegram.Firstname, this.Telegram.Lastname)
	}
	return this.Email
}

func (this User) GetAvatar(cdn string) string {
	if this.Telegram != nil && this.Telegram.Avatar != "" {
		return this.Telegram.Avatar
	}
	key := utils.Md5(this.Email)
	return fmt.Sprintf("%suser/avatar/%s", cdn, key)
}

func (this User) Balance(ctx context.Context, service *Service, config Config) (*big.Int, error) {
	query := `SELECT
		uw.user_id ,
		uw.salt,
		uw.wallet,
		SUM( IFNULL(d.eth, 0) ) AS deposit
	FROM
		adx.user_wallets AS uw
		LEFT JOIN adx.deposits AS d ON (d.user_id = uw.user_id)
	WHERE
		uw.user_id = %d
	AND uw.token_type = 'ETH'
	AND uw.is_main = 1
	GROUP BY
		user_id`
	db := service.Db
	rows, _, err := db.Query(query, this.Id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("not found")
	}
	row := rows[0]
	salt := row.Str(1)
	wallet := row.Str(2)
	deposit := new(big.Int).SetUint64(uint64(row.ForceFloat(3) * params.Shannon))
	deposit = new(big.Int).Mul(deposit, big.NewInt(params.Shannon))
	privateKey, err := utils.AddressDecrypt(wallet, salt, config.TokenSalt)
	if err != nil {
		return nil, err
	}
	publicKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	balance, err := service.Geth.BalanceAt(ctx, ethcommon.HexToAddress(publicKey), nil)
	if err != nil {
		return nil, err
	}
	balance = new(big.Int).Add(balance, deposit)
	if this.IsAdvertiser == 1 {
		query := `SELECT
			user_id ,
			SUM( price ) AS cost
		FROM
			adx.private_auctions
		WHERE
			audit_status < 2
		AND user_id = %d
		GROUP BY
			user_id`
		rows, _, err := db.Query(query, this.Id)
		if err != nil {
			return nil, err
		}
		if len(rows) > 0 {
			cost := new(big.Int).SetUint64(uint64(rows[0].ForceFloat(1) * params.Ether))
			balance = new(big.Int).Sub(balance, cost)
		}
	} else if this.IsPublisher == 1 {
		query := `SELECT
			a.user_id ,
			SUM( pa.price ) AS income
		FROM
			adx.private_auctions AS pa
		INNER JOIN adx.adzones AS a ON (a.id = pa.adzone_id)
		WHERE
			pa.audit_status = 1
		AND a.user_id = %d
		GROUP BY
			a.user_id`
		rows, _, err := db.Query(query, this.Id)
		if err != nil {
			return nil, err
		}
		if len(rows) > 0 {
			income := new(big.Int).SetUint64(uint64(rows[0].ForceFloat(1) * params.Ether))
			balance = new(big.Int).Add(balance, income)
		}
	}
	return balance, nil
}

func (this User) ETHBalance(ctx context.Context, service *Service, config Config) (*big.Int, error) {
	privateKey, err := utils.AddressDecrypt(this.Wallet, this.Salt, config.TokenSalt)
	if err != nil {
		return nil, err
	}
	publicKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return service.Geth.BalanceAt(ctx, ethcommon.HexToAddress(publicKey), nil)
}

func (this User) TokenBalance(ctx context.Context, service *Service, config Config, tokenAddress string) (*big.Int, error) {
	token, err := eth.NewToken(ethcommon.HexToAddress(tokenAddress), service.Geth)
	if err != nil {
		return nil, err
	}
	privateKey, err := utils.AddressDecrypt(this.Wallet, this.Salt, config.TokenSalt)
	if err != nil {
		return nil, err
	}
	publicKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return token.BalanceOf(nil, ethcommon.HexToAddress(publicKey))
}
