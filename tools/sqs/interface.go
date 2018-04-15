package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-msgqueue/msgqueue"
	"github.com/go-msgqueue/msgqueue/azsqs"
	"github.com/tokenme/adx/common"
)

type MsgType = uint

const (
	RegisterMsg MsgType = 0
	ResetPwdMsg MsgType = 1
)

type Message struct {
	Type         MsgType `codec:"type"`
	Email        string  `codec:"email"`
	Code         string  `codec:"code"`
	IsPublisher  uint    `code:"is_publisher"`
	IsAdvertiser uint    `code:"is_advertiser"`
}

func NewManager(config common.SQSConfig) msgqueue.Manager {
	return azsqs.NewManager(awsSQS(config), config.AccountId)
}

func NewQueue(config common.SQSConfig, opt *msgqueue.Options) msgqueue.Queue {
	m := NewManager(config)
	return m.NewQueue(opt)
}

func awsSQS(config common.SQSConfig) *sqs.SQS {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(config.EmailRegion),
		Credentials: credentials.NewStaticCredentials(config.AK, config.Secret, config.Token),
	}))
	return sqs.New(sess)
}
