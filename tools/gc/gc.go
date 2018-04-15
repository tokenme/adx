package gc

import (
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	"time"
)

const (
	VerifyCodeGCHours      int = 2
	AuthVerifyCodesGCHours int = 3
)

type Handler struct {
	Service *common.Service
	Config  common.Config
	exitCh  chan struct{}
}

func New(service *common.Service, config common.Config) *Handler {
	return &Handler{
		Service: service,
		Config:  config,
		exitCh:  make(chan struct{}, 1),
	}
}

func (this *Handler) Start() {
	log.Info("GC Start")
	hourlyTicker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-hourlyTicker.C:
			this.authVerifyCodesRecycle()
		case <-this.exitCh:
			hourlyTicker.Stop()
			return
		}
	}
}

func (this *Handler) Stop() {
	close(this.exitCh)
	log.Info("GC Stopped")
}

func (this *Handler) authVerifyCodesRecycle() error {
	db := this.Service.Db
	_, _, err := db.Query(`DELETE FROM adx.auth_verify_codes WHERE inserted<DATE_SUB(NOW(), INTERVAL %d HOUR)`, AuthVerifyCodesGCHours)
	return err
}
