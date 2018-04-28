package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-msgqueue/msgqueue"
	"github.com/go-msgqueue/msgqueue/azsqs"
	"github.com/tokenme/adx/common"
	"sync"
)

type MsgType = uint

const (
	RegisterMsg MsgType = 0
	ResetPwdMsg MsgType = 1
	AdClickMsg  MsgType = 2
	AdImpMsg    MsgType = 3
)

func NewManager(config common.SQSConfig) msgqueue.Manager {
	return azsqs.NewManager(awsSQS(config), config.AccountId)
}

func NewQueue(m msgqueue.Manager, opt *msgqueue.Options) msgqueue.Queue {
	return m.NewQueue(opt)
}

func awsSQS(config common.SQSConfig) *sqs.SQS {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AK, config.Secret, config.Token),
	}))
	return sqs.New(sess)
}

type AdQueue struct {
	ads []common.Ad
	sync.RWMutex
}

func (this *AdQueue) Add(ad common.Ad) {
	this.Lock()
	this.ads = append(this.ads, ad)
	this.Unlock()
}

func (this *AdQueue) Flush() (ads []common.Ad) {
	this.Lock()
	ads = this.ads
	this.ads = []common.Ad{}
	this.Unlock()
	return ads
}
