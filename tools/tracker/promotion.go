package tracker

import (
	"fmt"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	"strings"
	"sync"
	"time"
)

type PromotionLog struct {
	Proto        common.PromotionProto
	RecordOn     time.Time
	Pv           uint64
	Submissions  uint64
	Transactions uint64
	GiveOut      uint64
	Bonus        uint64
	ComissionFee uint64
}

type PromotionLogService struct {
	service *common.Service
	mp      map[string]*PromotionLog
	sync.Mutex
}

func NewPromotionLogService(service *common.Service) *PromotionLogService {
	return &PromotionLogService{
		service: service,
		mp:      make(map[string]*PromotionLog),
	}
}

func (this *PromotionLogService) Pv(proto common.PromotionProto) *PromotionLog {
	now := time.Now()
	ret := &PromotionLog{Proto: proto, RecordOn: now}
	trackKey := fmt.Sprintf("%d-%s", proto.Id, now.Format("2006-01-02"))
	this.Lock()
	if _, found := this.mp[trackKey]; !found {
		this.mp[trackKey] = ret
	}
	this.mp[trackKey].Pv += 1
	this.Unlock()
	return ret
}

func (this *PromotionLogService) Submission(proto common.PromotionProto) *PromotionLog {
	now := time.Now()
	ret := &PromotionLog{Proto: proto, RecordOn: now}
	trackKey := fmt.Sprintf("%d-%s", proto.Id, now.Format("2006-01-02"))
	this.Lock()
	if _, found := this.mp[trackKey]; !found {
		this.mp[trackKey] = ret
	}
	this.mp[trackKey].Submissions += 1
	this.Unlock()
	return ret
}

func (this *PromotionLogService) Transaction(proto common.PromotionProto) *PromotionLog {
	now := time.Now()
	ret := &PromotionLog{Proto: proto, RecordOn: now}
	trackKey := fmt.Sprintf("%d-%s", proto.Id, now.Format("2006-01-02"))
	this.Lock()
	if _, found := this.mp[trackKey]; !found {
		this.mp[trackKey] = ret
	}
	this.mp[trackKey].Transactions += 1
	this.Unlock()
	return ret
}

func (this *PromotionLogService) Transactions(proto common.PromotionProto, num uint64) *PromotionLog {
	now := time.Now()
	ret := &PromotionLog{Proto: proto, RecordOn: now}
	trackKey := fmt.Sprintf("%d-%s", proto.Id, now.Format("2006-01-02"))
	this.Lock()
	if _, found := this.mp[trackKey]; !found {
		this.mp[trackKey] = ret
	}
	this.mp[trackKey].Transactions += num
	this.Unlock()
	return ret
}

func (this *PromotionLogService) GiveOut(proto common.PromotionProto, giveOut uint64) *PromotionLog {
	now := time.Now()
	ret := &PromotionLog{Proto: proto, RecordOn: now}
	trackKey := fmt.Sprintf("%d-%s", proto.Id, now.Format("2006-01-02"))
	this.Lock()
	if _, found := this.mp[trackKey]; !found {
		this.mp[trackKey] = ret
	}
	this.mp[trackKey].GiveOut += giveOut
	this.Unlock()
	return ret
}

func (this *PromotionLogService) Bonus(proto common.PromotionProto, bonus uint64) *PromotionLog {
	now := time.Now()
	ret := &PromotionLog{Proto: proto, RecordOn: now}
	trackKey := fmt.Sprintf("%d-%s", proto.Id, now.Format("2006-01-02"))
	this.Lock()
	if _, found := this.mp[trackKey]; !found {
		this.mp[trackKey] = ret
	}
	this.mp[trackKey].Bonus += bonus
	this.Unlock()
	return ret
}

func (this *PromotionLogService) ComissionFee(proto common.PromotionProto, comissionFee uint64) *PromotionLog {
	now := time.Now()
	ret := &PromotionLog{Proto: proto, RecordOn: now}
	trackKey := fmt.Sprintf("%d-%s", proto.Id, now.Format("2006-01-02"))
	this.Lock()
	if _, found := this.mp[trackKey]; !found {
		this.mp[trackKey] = ret
	}
	this.mp[trackKey].ComissionFee += comissionFee
	this.Unlock()
	return ret
}

func (this *PromotionLogService) Flush() error {
	this.Lock()
	mp := this.mp
	this.mp = map[string]*PromotionLog{}
	this.Unlock()
	totalLogs := len(mp)
	if totalLogs == 0 {
		return nil
	}
	var val []string
	db := this.service.Db
	for _, v := range mp {
		val = append(val, fmt.Sprintf("(%d, '%s', %d, %d, %d, %d, %d, %d, %d, %d, %d, %d)", v.Proto.Id, v.RecordOn.Format("2006-01-02"), v.Proto.AdzoneId, v.Proto.ChannelId, v.Proto.UserId, v.Proto.AirdropId, v.Pv, v.Submissions, v.Transactions, v.GiveOut, v.Bonus, v.ComissionFee))
		if len(val) >= 1000 {
			_, _, err := db.Query(`INSERT INTO adx.promotion_stats (promotion_id, record_on, adzone_id, channel_id, promoter_id, airdrop_id, pv, submissions, transactions, give_out, bonus, commission_fee) VALUES %s ON DUPLICATE KEY UPDATE pv=VALUES(pv)+pv, submissions=VALUES(submissions)+submissions, transactions=VALUES(transactions)+transactions, give_out=VALUES(give_out)+give_out, bonus=VALUES(bonus)+bonus, commission_fee=VALUES(commission_fee)+commission_fee`, strings.Join(val, ","))
			if err != nil {
				log.Error(err.Error())
			}
			val = []string{}
		}
	}

	if len(val) > 0 {
		_, _, err := db.Query(`INSERT INTO adx.promotion_stats (promotion_id, record_on, adzone_id, channel_id, promoter_id, airdrop_id, pv, submissions, transactions, give_out, bonus, commission_fee) VALUES %s ON DUPLICATE KEY UPDATE pv=VALUES(pv)+pv, submissions=VALUES(submissions)+submissions, transactions=VALUES(transactions)+transactions, give_out=VALUES(give_out)+give_out, bonus=VALUES(bonus)+bonus, commission_fee=VALUES(commission_fee)+commission_fee`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
		}
		val = []string{}
	}
	return nil
}
