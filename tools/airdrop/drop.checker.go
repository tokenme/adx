package airdrop

import (
	"context"
	"fmt"
	"github.com/mkideal/log"
	ethutils "github.com/tokenme/adx/coins/eth/utils"
	"github.com/tokenme/adx/common"
	"github.com/tokenme/adx/tools/tracker"
	"gopkg.in/telegram-bot-api.v4"
	"time"
)

type AirdropChecker struct {
	service     *common.Service
	config      common.Config
	tracker     *tracker.Tracker
	telegramBot *tgbotapi.BotAPI
	exitCh      chan struct{}
}

func NewAirdropChecker(service *common.Service, config common.Config, tracker *tracker.Tracker) *AirdropChecker {
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Error(err.Error())
	}
	bot.Debug = config.Debug
	return &AirdropChecker{
		service:     service,
		config:      config,
		tracker:     tracker,
		telegramBot: bot,
		exitCh:      make(chan struct{}, 1),
	}
}

func (this *AirdropChecker) Start() {
	log.Info("AirdropChecker Start")
	ctx, cancel := context.WithCancel(context.Background())
	go this.CheckLoop(ctx)
	<-this.exitCh
	cancel()
}

func (this *AirdropChecker) Stop() {
	close(this.exitCh)
	log.Info("AirdropChecker Stopped")
}

func (this *AirdropChecker) CheckLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			this.Check(ctx)
		}
		time.Sleep(10 * time.Second)
	}
}

func (this *AirdropChecker) Check(ctx context.Context) {
	db := this.service.Db
	query := `SELECT tx FROM adx.airdrop_submissions WHERE tx>'%s' AND (status=1 OR tx_status=0) GROUP BY tx ORDER BY tx ASC LIMIT 1000`
	var (
		startTx string
		endTx   string
	)
	for {
		endTx = startTx
		rows, _, err := db.Query(query, startTx)
		if err != nil {
			log.Error(err.Error())
			break
		}
		for _, row := range rows {
			tx := row.Str(0)
			this.CheckSubmission(ctx, tx)
			endTx = tx
		}
		if endTx == startTx {
			break
		}
		startTx = endTx
	}
}

func (this *AirdropChecker) CheckSubmission(ctx context.Context, submissionTx string) error {
	receipt, err := ethutils.TransactionReceipt(this.service.Geth, ctx, submissionTx)
	if err != nil {
		//log.Error(err.Error())
		return err
	}
	if receipt == nil {
		log.Info("Submission Tx:%s, isPending", submissionTx)
		return nil
	}
	var (
		txStatus         uint = 2
		submissionStatus uint = 3
	)
	if receipt.Status == 1 {
		txStatus = 1
		submissionStatus = 2
	}
	log.Info("Submission Tx:%s, status:%d", submissionTx, txStatus)
	db := this.service.Db
	_, _, err = db.Query(`UPDATE adx.airdrop_tx SET status=%d WHERE tx='%s'`, txStatus, submissionTx)
	if err != nil {
		log.Error(err.Error())
	}
	_, _, err = db.Query(`UPDATE adx.airdrop_submissions SET status=%d, tx_status=%d WHERE tx='%s'`, submissionStatus, txStatus, submissionTx)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if txStatus != 1 {
		return nil
	}

	query := `SELECT
	ass.promotion_id,
	ass.adzone_id,
	ass.channel_id,
	ass.promoter_id,
	ass.airdrop_id,
	a.bonus,
	a.give_out,
	a.commission_fee,
	t.name,
	t.decimals,
	ass.telegram_msg_id,
	ass.telegram_chat_id
FROM adx.airdrop_submissions AS ass
INNER JOIN adx.airdrops AS a ON (a.id=ass.airdrop_id)
INNER JOIN adx.tokens AS t ON (t.address=a.token_address)
WHERE ass.tx='%s'`
	rows, _, err := db.Query(query, submissionTx)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	totalSubmissions := int64(len(rows))
	if totalSubmissions == 0 {
		return nil
	}
	row := rows[0]

	proto := common.PromotionProto{
		Id:        row.Uint64(0),
		AdzoneId:  row.Uint64(1),
		ChannelId: row.Uint64(2),
		UserId:    row.Uint64(3),
		AirdropId: row.Uint64(4),
	}
	airdrop := common.Airdrop{
		Id:            row.Uint64(4),
		Bonus:         row.Uint(5),
		GiveOut:       row.Uint64(6),
		CommissionFee: row.Uint64(7),
		Token: common.Token{
			Name:     row.Str(8),
			Decimals: row.Uint(9),
		},
	}

	this.tracker.Promotion.Transactions(proto, uint64(totalSubmissions))
	this.tracker.Promotion.GiveOut(proto, airdrop.TotalGiveOutDecimals(totalSubmissions).Uint64())
	this.tracker.Promotion.Bonus(proto, airdrop.TotalTokenBonusDecimals(totalSubmissions).Uint64())
	this.tracker.Promotion.ComissionFee(proto, airdrop.TotalCommissionFeeGwei(totalSubmissions).Uint64())

	if this.telegramBot != nil {
		for _, row := range rows {
			msgId := row.Int(10)
			chatId := row.Int64(11)
			if chatId != 0 {
				info := fmt.Sprintf("Congratulations! Your airdrop token (%d %s) has been send out!", airdrop.GiveOut, airdrop.Token.Name)
				msg := tgbotapi.NewMessage(chatId, info)
				msg.ReplyToMessageID = msgId
				this.telegramBot.Send(msg)
			}
		}
	}

	return nil
}
