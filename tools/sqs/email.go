package sqs

import (
	"bytes"
	"fmt"
	"github.com/go-msgqueue/msgqueue"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/url"
)

type EmailQueue struct {
	Service   *common.Service
	Config    common.Config
	Queue     msgqueue.Queue
	Processor *msgqueue.Processor
}

func NewEmailQueue(service *common.Service, config common.Config) *EmailQueue {
	queue := &EmailQueue{
		Service: service,
		Config:  config,
	}
	opt := &msgqueue.Options{
		Name:    config.SQS.EmailQueue,
		Handler: queue.Handler,
	}
	q := NewQueue(config.SQS, opt)
	queue.Queue = q
	queue.Processor = q.Processor()
	return queue
}

func (this *EmailQueue) Start() {
	this.Processor.Start()
}

func (this *EmailQueue) Stop() {
	this.Processor.Stop()
}

func (this *EmailQueue) NewRegister(user common.User) error {
	return this.Queue.Call(Message{Type: RegisterMsg, Email: user.Email, Code: user.ActivationCode, IsPublisher: user.IsPublisher, IsAdvertiser: user.IsAdvertiser})
}

func (this *EmailQueue) NewResetPwd(user common.User) error {
	return this.Queue.Call(Message{Type: ResetPwdMsg, Email: user.Email, Code: user.ResetPwdCode, IsPublisher: user.IsPublisher, IsAdvertiser: user.IsAdvertiser})
}

func (this *EmailQueue) Handler(msg Message) error {
	var (
		subject string
		body    string
	)
	switch msg.Type {
	case RegisterMsg:
		subject = "Welcome to Tokenama! Confirm Your Email"
		t, err := template.ParseFiles("templates/register-email.html")
		if err != nil {
			log.Error(err.Error())
			return err
		}
		var link string
		if msg.IsPublisher == 1 {
			link = this.Config.PublisherDomain
		} else {
			link = this.Config.AdvertiserDomain
		}
		query := &url.Values{}
		query.Add("email", msg.Email)
		query.Add("activation_code", msg.Code)
		query.Add("utm_campaign", "website")
		query.Add("utm_source", link)
		query.Add("utm_medium", "email")
		link = fmt.Sprintf("%s/user/activate?%s", link, query.Encode())
		var b bytes.Buffer
		t.Execute(&b, link)
		body = b.String()
	case ResetPwdMsg:
		subject = "Reset Password"
		t, err := template.ParseFiles("templates/reset-pwd-email.html")
		if err != nil {
			log.Error(err.Error())
			return err
		}
		var link string
		if msg.IsPublisher == 1 {
			link = this.Config.PublisherDomain
		} else {
			link = this.Config.AdvertiserDomain
		}
		query := &url.Values{}
		query.Add("email", msg.Email)
		query.Add("code", msg.Code)
		query.Add("utm_campaign", "website")
		query.Add("utm_source", link)
		query.Add("utm_medium", "email")
		link = fmt.Sprintf("%s/user/reset-pwd-verify?%s", link, query.Encode())
		var b bytes.Buffer
		t.Execute(&b, link)
		body = b.String()
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "support@tokenmama.io", "Tokenmama Support")
	m.SetHeader("To", msg.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewPlainDialer(this.Config.Mail.Server, this.Config.Mail.Port, this.Config.Mail.User, this.Config.Mail.Passwd)
	// Send the email to Bob, Cora and Dan.
	err := d.DialAndSend(m)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
