package telegram

import (
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/tokenme/adx/common"
	"github.com/tokenme/adx/tools/shorturl"
	"github.com/tokenme/adx/tools/tracker"
	"github.com/tokenme/adx/utils/token"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
	"strconv"
	"sync"
)

type Bot struct {
	Service     *common.Service
	Config      common.Config
	TelegramBot *tgbotapi.BotAPI
	UpdatesCh   tgbotapi.UpdatesChannel
	Tracker     *tracker.Tracker
	wg          *sync.WaitGroup
}

type Message struct {
	Message *tgbotapi.Message
	Code    token.Token
}

func New(service *common.Service, config common.Config, trackerService *tracker.Tracker) (*Bot, error) {
	if !config.EnableTelegramBot {
		return nil, errors.New("telemebot disabled")
	}
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	bot.Debug = config.Debug
	bot.Buffer = 1000
	return &Bot{
		Service:     service,
		Config:      config,
		TelegramBot: bot,
		Tracker:     trackerService,
		wg:          &sync.WaitGroup{},
	}, nil
}

func (this *Bot) Start() error {
	if this.TelegramBot == nil {
		return errors.New("missing telegram bot")
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updateConfig.Limit = 100
	updates, err := this.TelegramBot.GetUpdatesChan(updateConfig)
	if err != nil {
		return err
	}
	this.UpdatesCh = updates
	messageBuf := make(chan Message, 1000)
	pool, _ := ants.NewPoolWithFunc(500, func(msg interface{}) error {
		defer this.wg.Done()
		this.VerifyCodeHandler(msg.(Message))
		return nil
	})
	go func() {
		for {
			select {
			case msg := <-messageBuf:
				this.wg.Add(1)
				pool.Serve(msg)
			}
		}
	}()
	go func(updates tgbotapi.UpdatesChannel) {
		for update := range updates {
			this.MessageConsumer(update.Message, messageBuf)
		}
	}(updates)
	log.Info("TelegramBot Started")
	return nil
}

func (this *Bot) Stop() {
	log.Info("TelegramBot Stopped")
	if this.TelegramBot == nil {
		return
	}
	this.wg.Wait()
	this.UpdatesCh.Clear()
}

func (this *Bot) MessageConsumer(message *tgbotapi.Message, messageBuf chan Message) {
	if message == nil {
		return
	}
	if message.From == nil || message.Chat == nil {
		return
	}

	reg := regexp.MustCompile(fmt.Sprintf(`^/(\w+)@%s`, this.Config.TelegramBotName))
	matches := reg.FindStringSubmatch(message.Text)
	if len(matches) != 2 {
		return
	}
	cmd := matches[1]
	verifyCode, err := token.Decode(cmd)
	if err != nil {
		reply := "Sorry, submitted verification code is invalid"
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		msg.ReplyToMessageID = message.MessageID
		this.TelegramBot.Send(msg)
		return
	}
	log.Info("Get chat:%s, code:%s", message.Chat.UserName, cmd)

	messageBuf <- Message{Message: message, Code: verifyCode}
}

func (this *Bot) VerifyCodeHandler(msg Message) {
	var reply, lang string
	message := msg.Message
	log.Info("Verify chat:%s, user:%s, code:%d", message.Chat.UserName, message.From.UserName, msg.Code)
	db := this.Service.Db
	rows, _, err := db.Query(`SELECT c.status, c.promotion_id, c.adzone_id, c.channel_id, c.promoter_id, c.airdrop_id, a.telegram_group, c.wallet, c.referrer, c.email, a.telegram_admin, a.reply_msg FROM tokenme.codes AS c LEFT JOIN tokenme.airdrops AS a ON (a.id=c.airdrop_id) WHERE c.id=%d LIMIT 1`, msg.Code)
	if err != nil {
		reply = il8n_trans("Sorry, we have some internal server bug :(", lang)
	}
	if len(rows) == 0 {
		reply = il8n_trans("Sorry, submitted verification code is invalid or expired", lang)
	} else {
		codeStatus := rows[0].Uint(0)
		proto := common.PromotionProto{
			Id:        rows[0].Uint64(1),
			AdzoneId:  rows[0].Uint64(2),
			ChannelId: rows[0].Uint64(3),
			UserId:    rows[0].Uint64(4),
			AirdropId: rows[0].Uint64(5),
		}
		if proto.AirdropId == 51 { // 火币 空头ID
			lang = "vi-VN" // 越南语
		}
		telegramGroupName := rows[0].Str(6)
		wallet := rows[0].Str(7)
		referrer := rows[0].Str(8)
		email := rows[0].Str(9)
		telegramAdmin := rows[0].Str(10)
		replyMsg := rows[0].Str(11)
		if email == "" {
			email = "NULL"
		} else {
			email = fmt.Sprintf("'%s'", db.Escape(email))
		}
		if telegramGroupName != message.Chat.UserName {
			reply = fmt.Sprintf(il8n_trans("Sorry, you must submit your code in group", lang)+" @%s", telegramGroupName)
		} else if codeStatus == 0 {
			reply = il8n_trans("Sorry, maybe you didn't submit your wallet address?", lang)
		} else if codeStatus == 1 {
			telegramUsername := db.Escape(message.From.UserName)
			telegramUserId := strconv.FormatInt(int64(message.From.ID), 10)
			if telegramUsername != "" && telegramUsername == telegramAdmin {
				telegramUserId = "NULL"
			}
			_, ret, err := db.Query(`INSERT IGNORE INTO tokenme.airdrop_submissions (promotion_id, adzone_id, channel_id, promoter_id, airdrop_id, verify_code, email, wallet, telegram_msg_id, telegram_chat_id, telegram_user_id, telegram_chat_title, telegram_username, telegram_user_firstname, telegram_user_lastname, referrer) VALUES (%d, %d, %d, %d, %d, %d, %s, '%s', %d, %d, %s, '%s', '%s', '%s', '%s', '%s')`, proto.Id, proto.AdzoneId, proto.ChannelId, proto.UserId, proto.AirdropId, msg.Code, email, db.Escape(wallet), message.MessageID, message.Chat.ID, telegramUserId, db.Escape(message.Chat.UserName), telegramUsername, db.Escape(message.From.FirstName), db.Escape(message.From.LastName), db.Escape(referrer))
			if err != nil {
				log.Error(err.Error())
				reply = il8n_trans("Sorry, we have some internal server bug :(", lang)
			}
			_, _, err = db.Query(`UPDATE tokenme.codes SET status=2 WHERE id=%d`, msg.Code)
			if err != nil {
				reply = il8n_trans("Sorry, we have some internal server bug :(", lang)
			} else if ret.AffectedRows() == 0 {
				reply = il8n_trans("Sorry, you already submitted in this airdrop and could not submit again", lang)
			} else {
				if replyMsg == "" {
					reply = il8n_trans("Great! please wait for the airdrop transaction complete.", lang)
				} else {
					reply = replyMsg
				}
				proto.Referrer = wallet
				promoKey, err := common.EncodePromotion([]byte(this.Config.LinkSalt), proto)
				if err == nil {
					link := fmt.Sprintf("%s/promo/%s", this.Config.BaseUrl, promoKey)
					shortURL, err := shorturl.Sina(link)
					if err == nil && shortURL != "" {
						link = shortURL
					}
					reply = fmt.Sprintf("%s %s", reply, link)
				}
				this.Tracker.Promotion.Submission(proto)
			}
		} else if codeStatus == 2 {
			reply = il8n_trans("Your airdrop is in the pool, please wait for the transaction complete", lang)
		}
	}
	replyMsg := tgbotapi.NewMessage(message.Chat.ID, reply)
	replyMsg.ReplyToMessageID = message.MessageID
	this.TelegramBot.Send(replyMsg)
}
