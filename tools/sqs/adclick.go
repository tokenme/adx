package sqs

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/go-msgqueue/msgqueue"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	"time"
)

type AdClickQueue struct {
	Service   *common.Service
	Config    common.Config
	Queue     msgqueue.Queue
	Processor *msgqueue.Processor
	AdQueue   *AdQueue
	exitCh    chan struct{}
}

func NewAdClickQueue(m msgqueue.Manager, service *common.Service, config common.Config) *AdClickQueue {
	queue := &AdClickQueue{
		Service: service,
		Config:  config,
		AdQueue: &AdQueue{},
		exitCh:  make(chan struct{}, 1),
	}
	opt := &msgqueue.Options{
		Name:    config.SQS.AdClickQueue,
		Handler: queue.Handler,
	}
	q := m.NewQueue(opt)
	queue.Queue = q
	queue.Processor = q.Processor()
	return queue
}

func (this *AdClickQueue) Start() {
	this.Processor.Start()
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				this.Flush()
			case <-this.exitCh:
				ticker.Stop()
				this.Flush()
				return
			}
		}
	}()
}

func (this *AdClickQueue) Stop() {
	this.Processor.Stop()
	close(this.exitCh)
}

func (this *AdClickQueue) NewClick(msg string) error {
	return this.Queue.Call(msg)
}

func (this *AdClickQueue) Handler(msg string) error {
	ad, err := common.DecodeAd([]byte(this.Config.LinkSalt), msg)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	this.AdQueue.Add(ad)
	return nil
}

func (this *AdClickQueue) Flush() error {
	ads := this.AdQueue.Flush()
	if len(ads) == 0 {
		return nil
	}
	tx, err := this.Service.Clickhouse.Begin()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO adx.clicks (LogDate, LogTime, ReqId, CreativeId, AuctionId, AdzoneId, MediaId, SizeId, AdvertiserId, PublisherId, IP, Cookie, Link, Referrer, ScreenSize, AdSize, OsName, OsVersion, BrowserName, BrowserVersion, BrowserType, CountryId, CountryName, UserAgent) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	for _, ad := range ads {
		logTime := time.Unix(ad.LogTime, 0)
		if _, err := stmt.Exec(
			logTime,
			logTime,
			ad.PK,
			ad.Id,
			ad.AuctionId,
			ad.AdzoneId,
			ad.MediaId,
			ad.SizeId,
			ad.AdvertiserId,
			ad.PublisherId,
			ad.Env.IP,
			ad.Env.Cookie,
			ad.Env.URL,
			ad.Env.Referrer,
			ad.Env.ScreenSize,
			ad.Env.AdSize,
			ad.Env.OsName,
			ad.Env.OsVersion,
			ad.Env.BrowserName,
			ad.Env.BrowserVersion,
			uint16(ad.Env.BrowserType),
			uint32(ad.Env.CountryId),
			ad.Env.CountryName,
			ad.Env.UserAgent,
		); err != nil {
			log.Error(err.Error())
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
