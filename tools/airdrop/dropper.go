package airdrop

import (
	"context"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/coins/eth"
	ethutils "github.com/tokenme/adx/coins/eth/utils"
	"github.com/tokenme/adx/common"
	"github.com/tokenme/adx/utils"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Airdropper struct {
	service *common.Service
	config  common.Config
	wg      *sync.WaitGroup
	stopCh  chan struct{}
	exitCh  chan struct{}
}

func NewAirdropper(service *common.Service, config common.Config) *Airdropper {
	return &Airdropper{
		service: service,
		config:  config,
		exitCh:  make(chan struct{}, 1),
		stopCh:  make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},
	}
}

func (this *Airdropper) Start() {
	log.Info("Airdropper Start")
	ctx, cancel := context.WithCancel(context.Background())
	go this.DropLoop(ctx)
	<-this.exitCh
	cancel()
}

func (this *Airdropper) Stop() {
	this.stopCh <- struct{}{}
	this.wg.Wait()
	close(this.exitCh)
	log.Info("Airdropper Stopped")
}

func (this *Airdropper) DropLoop(ctx context.Context) {
	var interval = time.Duration(5)
	newDrop := this.Drop(ctx)
	if newDrop {
		interval = time.Duration(1)
	} else {
		interval = time.Duration(5)
	}
	time.Sleep(interval * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-this.stopCh:
			return
		default:
			newDrop := this.Drop(ctx)
			if newDrop {
				interval = time.Duration(1)
			} else {
				interval = time.Duration(5)
			}
		}
		time.Sleep(interval * time.Minute)
	}
}

func (this *Airdropper) Drop(ctx context.Context) bool {
	db := this.service.Db
	query := `SELECT
	a.id ,
	a.wallet ,
	a.salt ,
	a.token_address ,
	a.gas_price ,
	a.gas_limit ,
	a.bonus ,
	a.commission_fee ,
	a.give_out ,
	t.decimals ,
	a.dealer_contract,
	a.sync_drop
FROM
	adx.airdrops AS a
INNER JOIN adx.tokens AS t ON ( t.address = a.token_address )
WHERE
	t.protocol = 'ERC20'
AND a.balance_status = 0
AND ((a.approve_tx_status = 2 AND a.dealer_tx_status = 2) OR a.sync_drop=1)
AND EXISTS ( SELECT
	1
FROM
	adx.airdrop_submissions AS ass
WHERE
	ass.status IN (0, 2)
AND ass.airdrop_id = a.id
AND ass.blocked = 0
LIMIT 1 )
AND a.id > %d
AND a.drop_date <= DATE( NOW())
ORDER BY
	a.id DESC
LIMIT 100`
	var (
		startId uint64
		endId   uint64
		ret     bool
	)
	for {
		endId = startId
		rows, _, err := db.Query(query, startId)
		if err != nil {
			log.Error(err.Error())
			break
		}
		var airdrops []*common.Airdrop
		var wg sync.WaitGroup
		for _, row := range rows {
			wallet := row.Str(1)
			salt := row.Str(2)
			privateKey, _ := utils.AddressDecrypt(wallet, salt, this.config.TokenSalt)
			publicKey, _ := eth.AddressFromHexPrivateKey(privateKey)
			airdrop := &common.Airdrop{
				Id:             row.Uint64(0),
				Wallet:         publicKey,
				WalletPrivKey:  privateKey,
				Token:          common.Token{Address: row.Str(3), Decimals: row.Uint(9)},
				GasPrice:       row.Uint64(4),
				GasLimit:       row.Uint64(5),
				Bonus:          row.Uint(6),
				CommissionFee:  row.Uint64(7),
				GiveOut:        row.Uint64(8),
				DealerContract: row.Str(10),
				SyncDrop:       row.Uint(11),
			}
			endId = airdrop.Id
			wg.Add(1)
			go func(airdrop *common.Airdrop, c context.Context) {
				defer wg.Done()
				airdrop.CheckBalance(this.service.Geth, c)
			}(airdrop, ctx)
			airdrops = append(airdrops, airdrop)
		}
		wg.Wait()
		if len(airdrops) > 0 {
			ret = true
		}
		for _, airdrop := range airdrops {
			this.DropAirdrop(ctx, airdrop)
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	return ret
}

func (this *Airdropper) DropAirdrop(ctx context.Context, airdrop *common.Airdrop) {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT COUNT(*) AS num FROM adx.airdrop_submissions WHERE status IN (0, 3) AND airdrop_id=%d`, airdrop.Id)
	if err != nil {
		log.Error(err.Error())
		return
	}
	totalSubmissions := rows[0].Int64(0)
	if totalSubmissions == 0 {
		return
	}
	/*
			gasNeed, tokenNeed, enoughGas, enoughToken := airdrop.EnoughBudgetForSubmissions(totalSubmissions)
			if !enoughGas {
				log.Error("Not enough gas, need:%d, left:%d", gasNeed.Uint64(), airdrop.GasBalance.Uint64())
				return
			}
			if !enoughToken {
				log.Error("Not enough token, need:%d, left:%d", tokenNeed.Uint64(), airdrop.TokenBalance.Uint64())
				return
			}

		if airdrop.SyncDrop == 0 {
			token, err := ethutils.NewStandardToken(airdrop.Token.Address, this.service.Geth)
			if err != nil {
				log.Error(err.Error())
				return
			}
			allowance, err := ethutils.StandardTokenAllowance(token, airdrop.Wallet, airdrop.DealerContract)
			if err != nil {
				log.Error(err.Error())
				return
			}
			if allowance.Cmp(tokenNeed) == -1 {
				db.Query("UPDATE tokenme.airdrops SET allowance=0, approve_tx_status=0 WHERE id=%d", airdrop.Id)
				log.Error("Not enough allowance, need:%d, left:%d, contract:%s", tokenNeed.Uint64(), allowance.Uint64(), airdrop.DealerContract)
				return
			}
		}
	*/
	query := `SELECT
	ass.id ,
	ass.promotion_id ,
	ass.adzone_id ,
	ass.channel_id ,
	ass.promoter_id ,
	ass.wallet ,
	ass.referrer ,
	u.wallet ,
	u.salt
FROM
	adx.airdrop_submissions AS ass
INNER JOIN adx.user_wallets AS u ON ( u.user_id = ass.promoter_id
AND u.token_type = 'ETH'
AND u.is_main = 1 )
WHERE
	ass.status IN (0, 3)
AND ass.blocked = 0
AND ass.airdrop_id = %d
ORDER BY
	id DESC LIMIT 1000`
	var submissions []*common.AirdropSubmission
	rows, _, err = db.Query(query, airdrop.Id)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, row := range rows {
		wallet := row.Str(7)
		salt := row.Str(8)
		privateKey, _ := utils.AddressDecrypt(wallet, salt, this.config.TokenSalt)
		publicKey, _ := eth.AddressFromHexPrivateKey(privateKey)
		submissionWallet := row.Str(5)
		referrer := row.Str(6)
		if referrer == submissionWallet || referrer == publicKey {
			referrer = ""
		}
		submission := &common.AirdropSubmission{
			Id:      row.Uint64(0),
			Airdrop: airdrop,
			Proto: common.PromotionProto{
				Id:        row.Uint64(1),
				AdzoneId:  row.Uint64(2),
				ChannelId: row.Uint64(3),
				UserId:    row.Uint64(4),
				Referrer:  referrer,
			},
			Wallet:         submissionWallet,
			PromoterWallet: publicKey,
		}
		submissions = append(submissions, submission)
	}
	if airdrop.SyncDrop == 0 {
		this.DropAirdropChunk(ctx, airdrop, submissions)
	} else {
		this.DropAirdropSync(ctx, airdrop, submissions)
	}
}

func (this *Airdropper) PrepareAirdrop(ctx context.Context, airdrop *common.Airdrop, submissions []*common.AirdropSubmission) {

	totalSubmissions := int64(len(submissions))
	if totalSubmissions == 0 {
		return
	}
	log.Info("Prepareing %d submissions for airdrop:%d", totalSubmissions, airdrop.Id)
	var (
		recipientsMap = make(map[string]*big.Int)
		promoWallet   = submissions[0].PromoterWallet
	)
	for _, submission := range submissions {
		recipientsMap[submission.Wallet] = airdrop.TokenGiveOut()
		if submission.Proto.Referrer != "" {
			if _, found := recipientsMap[submission.Proto.Referrer]; found {
				recipientsMap[submission.Proto.Referrer] = new(big.Int).Add(recipientsMap[submission.Proto.Referrer], airdrop.TokenBonus())
			} else {
				recipientsMap[submission.Proto.Referrer] = airdrop.TokenBonus()
			}
		} else if promoWallet != submission.Wallet {
			if _, found := recipientsMap[promoWallet]; found {
				recipientsMap[promoWallet] = new(big.Int).Add(recipientsMap[promoWallet], airdrop.TokenBonus())
			} else {
				recipientsMap[promoWallet] = airdrop.TokenBonus()
			}
		}
	}
	var val []string
	db := this.service.Db
	decimalsPow := new(big.Int).SetUint64(uint64(math.Pow10(int(airdrop.Token.Decimals))))
	for addr, amount := range recipientsMap {
		value := new(big.Int).Div(amount, decimalsPow)
		val = append(val, fmt.Sprintf("(%d, '%s', %d)", airdrop.Id, db.Escape(addr), value.Uint64()))
		if len(val) > 1000 {
			_, _, err := db.Query(`INSERT INTO adx.airdrop_prepares (airdrop_id, wallet, amount) VALUES %s ON DUPLICATE KEY UPDATE amount=VALUES(amount)`, strings.Join(val, ","))
			if err != nil {
				log.Error(err.Error())
			}
			val = []string{}
		}
	}

	if len(val) > 0 {
		_, _, err := db.Query(`INSERT INTO adx.airdrop_prepares (airdrop_id, wallet, amount) VALUES %s ON DUPLICATE KEY UPDATE amount=VALUES(amount)`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
		}
		val = []string{}
	}
}

func (this *Airdropper) DropAirdropSync(ctx context.Context, airdrop *common.Airdrop, submissions []*common.AirdropSubmission) {
	totalSubmissions := int64(len(submissions))
	if totalSubmissions == 0 {
		return
	}
	log.Info("Sync Airdroping %d submissions for airdrop:%d", totalSubmissions, airdrop.Id)
	var (
		promoWallet     = submissions[0].PromoterWallet
		promotionAmount *big.Int
	)
	for _, submission := range submissions {
		amount := airdrop.TokenGiveOut()
		addr := submission.Wallet
		this.transfer(ctx, airdrop, submission.Id, addr, amount, false)
		if submission.Proto.Referrer != "" {
			amount := airdrop.TokenBonus()
			addr := submission.Proto.Referrer
			this.transfer(ctx, airdrop, submission.Id, addr, amount, true)
		} else if promoWallet != submission.Wallet {
			if promotionAmount != nil {
				promotionAmount = new(big.Int).Add(promotionAmount, airdrop.TokenBonus())
			} else {
				promotionAmount = airdrop.TokenBonus()
			}
		}
	}
	if promotionAmount != nil {
		this.transfer(ctx, airdrop, 0, promoWallet, promotionAmount, true)
	}
}

func (this *Airdropper) transfer(ctx context.Context, airdrop *common.Airdrop, submissionId uint64, addr string, amount *big.Int, isReferrer bool) error {
	log.Info("W: %s, V: %d", addr, amount)
	token, err := ethutils.NewToken(airdrop.Token.Address, this.service.Geth)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	transactor := eth.TransactorAccount(airdrop.WalletPrivKey)
	nonce, err := eth.Nonce(ctx, this.service.Geth, this.service.Redis.Master, airdrop.Wallet, "main")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: airdrop.GasPriceToWei(),
		GasLimit: airdrop.GasLimit,
	}
	eth.TransactorUpdate(transactor, transactorOpts, ctx)
	tx, err := ethutils.Transfer(token, transactor, addr, amount)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	eth.NonceIncr(ctx, this.service.Geth, this.service.Redis.Master, airdrop.Wallet, "main")
	if !isReferrer {
		txHash := tx.Hash()
		db := this.service.Db
		_, _, err = db.Query(`UPDATE adx.airdrop_submissions SET status=1, tx='%s' WHERE airdrop_id=%d AND id=%d`, db.Escape(txHash.Hex()), airdrop.Id, submissionId)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (this *Airdropper) DropAirdropChunk(ctx context.Context, airdrop *common.Airdrop, submissions []*common.AirdropSubmission) {
	totalSubmissions := int64(len(submissions))
	if totalSubmissions == 0 {
		return
	}
	log.Info("Async Airdroping %d submissions for airdrop:%d", totalSubmissions, airdrop.Id)
	var (
		recipientsMap = make(map[string]*big.Int)
		promoWallet   = submissions[0].PromoterWallet
		submissionMap = make(map[string]uint64)
	)
	for _, submission := range submissions {
		recipientsMap[submission.Wallet] = airdrop.TokenGiveOut()
		submissionMap[submission.Wallet] = submission.Id
		if submission.Proto.Referrer != "" {
			if _, found := recipientsMap[submission.Proto.Referrer]; found {
				recipientsMap[submission.Proto.Referrer] = new(big.Int).Add(recipientsMap[submission.Proto.Referrer], airdrop.TokenBonus())
			} else {
				recipientsMap[submission.Proto.Referrer] = airdrop.TokenBonus()
			}
		} else if promoWallet != submission.Wallet {
			if _, found := recipientsMap[promoWallet]; found {
				recipientsMap[promoWallet] = new(big.Int).Add(recipientsMap[promoWallet], airdrop.TokenBonus())
			} else {
				recipientsMap[promoWallet] = airdrop.TokenBonus()
			}
		}

	}

	var (
		tokenAmounts  []*big.Int
		recipients    []ethcommon.Address
		submissionIds []uint64
	)
	for addr, amount := range recipientsMap {
		recipients = append(recipients, ethcommon.HexToAddress(addr))
		tokenAmounts = append(tokenAmounts, amount)
		if id, found := submissionMap[addr]; found {
			submissionIds = append(submissionIds, id)
		}
		if len(recipients) >= 10 {
			this.DropChunk(ctx, airdrop, submissionIds, recipients, tokenAmounts)
			recipients = []ethcommon.Address{}
			tokenAmounts = []*big.Int{}
			submissionIds = []uint64{}
		}
	}
	if len(recipients) > 0 {
		this.DropChunk(ctx, airdrop, submissionIds, recipients, tokenAmounts)
		recipients = []ethcommon.Address{}
		tokenAmounts = []*big.Int{}
		submissionIds = []uint64{}
	}

	log.Info("Async Done %d submissions for airdrop:%d", totalSubmissions, airdrop.Id)
}

func (this *Airdropper) DropChunk(ctx context.Context, airdrop *common.Airdrop, submissionIds []uint64, recipients []ethcommon.Address, tokenAmounts []*big.Int) error {
	this.wg.Add(1)
	defer this.wg.Done()
	airdropper, err := ethutils.NewAirdropper(airdrop.DealerContract, this.service.Geth)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	transactor := eth.TransactorAccount(airdrop.WalletPrivKey)

	nonce, err := eth.Nonce(ctx, this.service.Geth, this.service.Redis.Master, airdrop.Wallet, "main")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: airdrop.GasPriceToWei(),
		GasLimit: this.config.Airdrop.DealerContractGasLimit,
	}
	eth.TransactorUpdate(transactor, transactorOpts, ctx)
	tx, err := ethutils.AirdropperDrop(airdropper, transactor, recipients, tokenAmounts)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	eth.NonceIncr(ctx, this.service.Geth, this.service.Redis.Master, airdrop.Wallet, "main")
	txHash := tx.Hash()
	log.Info("tx: %s, nonce: %d", txHash.Hex(), nonce)
	db := this.service.Db
	_, _, err = db.Query(`INSERT IGNORE INTO adx.airdrop_tx (tx, airdrop_id) VALUES ('%s', %d)`, txHash.Hex(), airdrop.Id)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if len(submissionIds) > 0 {
		var ids []string
		for _, id := range submissionIds {
			ids = append(ids, strconv.FormatUint(id, 10))
		}
		_, _, err = db.Query(`UPDATE adx.airdrop_submissions SET status=1, tx='%s' WHERE airdrop_id=%d AND id IN (%s)`, db.Escape(txHash.Hex()), airdrop.Id, strings.Join(ids, ","))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}
