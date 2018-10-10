package common

type Channel struct {
	Id     uint64 `json:"id,omitempty"`
	UserId uint64 `json:"user_id,omitempty"`
	Name   string `json:"name,omitempty"`
}
