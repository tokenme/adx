package tracker

import (
	"github.com/mkideal/log"
	"time"
)

type Tracker struct {
	Promotion *PromotionLogService
	exitCh    chan struct{}
}

func New() *Tracker {
	return &Tracker{
		exitCh: make(chan struct{}, 1),
	}
}

func (this *Tracker) SetPromotion(promotion *PromotionLogService) {
	this.Promotion = promotion
}

func (this *Tracker) Start() {
	log.Info("Tracker Start")
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			this.Flush()
		case <-this.exitCh:
			ticker.Stop()
			return
		}
	}
}

func (this *Tracker) Stop() {
	close(this.exitCh)
	log.Info("Tracker Stopped")
}

func (this *Tracker) Flush() {
	if this.Promotion != nil {
		this.Promotion.Flush()
	}
}
