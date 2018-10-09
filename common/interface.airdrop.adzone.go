package common

type AirdropAdzone struct {
	Id        uint64 `json:"id,omitempty"`
	UserId    uint64 `json:"user_id,omitempty"`
	ChannelId uint64 `json:"channel_id,omitempty"`
	Name      string `json:"name,omitempty"`
}
