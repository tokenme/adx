package common

import (
	"time"
)

type Settlement = uint

const (
	CPM  Settlement = 1
	CPT  Settlement = 2
	CPMT Settlement = 3
)

type Adzone struct {
	Id           uint64     `json:"id"`
	Media        Media      `json:"media"`
	Size         Size       `json:"size"`
	Url          string     `json:"url"`
	MinCPT       float64    `json:"min_cpt,omitempty"`
	MinCPM       float64    `json:"min_cpm,omitempty"`
	Settlement   Settlement `json:"settlement"`
	Desc         string     `json:"desc"`
	Rolling      uint       `json:"rolling"`
	OnlineStatus uint       `json:"online_status"`
	InsertedAt   time.Time  `json:"inserted_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
