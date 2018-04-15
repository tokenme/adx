package common

import (
	"fmt"
	"github.com/tokenme/adx/utils"
)

type User struct {
	Id             uint64        `json:"id,omitempty"`
	Email          string        `json:"email,omitempty"`
	ActivationCode string        `json:"-"`
	ResetPwdCode   string        `json:"-"`
	Telegram       *TelegramUser `json:"telegram,omitempty"`
	ShowName       string        `json:"showname,omitempty"`
	Avatar         string        `json:"avatar,omitempty"`
	Salt           string        `json:"-"`
	Password       string        `json:"-"`
	IsAdmin        uint          `json:"is_admin,omitempty"`
	IsPublisher    uint          `json:"is_publisher,omitempty"`
	IsAdvertiser   uint          `json:"is_advertiser,omitempty"`
}

type TelegramUser struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Avatar    string `json:"avatar"`
}

func (this User) GetShowName() string {
	if this.Telegram != nil && (this.Telegram.Firstname != "" && this.Telegram.Lastname != "" || this.Telegram.Username != "") {
		if this.Telegram.Username != "" {
			return this.Telegram.Username
		}
		return fmt.Sprintf("%s %s", this.Telegram.Firstname, this.Telegram.Lastname)
	}
	return this.Email
}

func (this User) GetAvatar(cdn string) string {
	if this.Telegram != nil && this.Telegram.Avatar != "" {
		return this.Telegram.Avatar
	}
	key := utils.Md5(this.Email)
	return fmt.Sprintf("%suser/avatar/%s", cdn, key)
}
