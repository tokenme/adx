package common

import (
	"time"
)

type PrivateAuction struct {
	Id           uint64                   `json:"id"`
	UserId       uint64                   `json:"user_id,omitempty"`
	Adzone       Adzone                   `json:"adzone,omitempty"`
	Title        string                   `json:"title,omitempty"`
	Price        float64                  `json:"price,omitempty"`
	Cost         float64                  `json:"cost,omitempty"`
	StartTime    time.Time                `json:"start_time"`
	EndTime      time.Time                `json:"end_time"`
	AuditStatus  uint                     `json:"audit_status"`
	OnlineStatus uint                     `json:"online_status"`
	RejectReason string                   `json:"reject_reason,omitempty"`
	Creatives    []PrivateAuctionCreative `json:"creatives,omitempty"`
	InsertedAt   time.Time                `json:"inserted_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

func (this PrivateAuction) GetCost() float64 {
	days := this.EndTime.Sub(this.StartTime).Hours()/24 + 1
	return this.Price * 100000 * days / 100000
}
