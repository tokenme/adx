package server

import (
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	"time"
)

type Server struct {
	service             *common.Service
	config              common.Config
	privateAuctionsPool *PrivateAuctionPool
	exitCh              chan struct{}
}

func New(service *common.Service, config common.Config) *Server {
	return &Server{
		service:             service,
		config:              config,
		privateAuctionsPool: NewPrivateAuctionPool(),
		exitCh:              make(chan struct{}, 1),
	}
}

func (this *Server) Start() {
	log.Info("AdServer Start")
	ticker := time.NewTicker(1 * time.Minute)
	this.privateAuctionsPool.Reload(this.service, this.config)
	for {
		select {
		case <-ticker.C:
			this.privateAuctionsPool.Reload(this.service, this.config)
		case <-this.exitCh:
			ticker.Stop()
			return
		}
	}
}

func (this *Server) Stop() {
	close(this.exitCh)
	log.Info("AdServer Stopped")
}

func (this *Server) PrivateAuction(adzoneId uint64) (ad common.Ad, err error) {
	return this.privateAuctionsPool.Pop(adzoneId)
}
