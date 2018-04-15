package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/tokenme/adx/common"
	"github.com/tokenme/adx/utils"
	"sort"
	"strings"
	"time"
)

func ParseTelegramAuth(data string) (telegram common.TelegramUser, err error) {
	t := make(map[string]interface{})
	err = json.Unmarshal([]byte(data), &t)
	if err != nil {
		return telegram, err
	}
	var telegramId int64
	switch t["id"].(type) {
	case float64:
		telegramId = int64(t["id"].(float64))
	case int:
		telegramId = int64(t["id"].(int))
	case int64:
		telegramId = t["id"].(int64)
	case uint64:
		telegramId = int64(t["id"].(uint64))
	}
	var (
		username  string
		firstname string
		lastname  string
		photoUrl  string
	)
	if v, found := t["username"]; found {
		username = v.(string)
	}
	if v, found := t["first_name"]; found {
		firstname = v.(string)
	}
	if v, found := t["last_name"]; found {
		lastname = v.(string)
	}
	if v, found := t["photo_url"]; found {
		photoUrl = v.(string)
	}
	return common.TelegramUser{
		Id:        telegramId,
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Avatar:    photoUrl,
	}, nil
}

func TelegramAuthCheck(data string, botToken string) bool {
	telegram := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &telegram)
	if err != nil {
		return false
	}
	var (
		checkArr  []string
		checkHash = telegram["hash"].(string)
	)
	for key, val := range telegram {
		if key == "hash" {
			continue
		}
		switch val.(type) {
		case string:
			checkArr = append(checkArr, fmt.Sprintf("%s=%s", key, val.(string)))
		case float64:
			checkArr = append(checkArr, fmt.Sprintf("%s=%d", key, int64(val.(float64))))
		case int:
			checkArr = append(checkArr, fmt.Sprintf("%s=%d", key, int64(val.(int))))
		case int64:
			checkArr = append(checkArr, fmt.Sprintf("%s=%d", key, val.(int64)))
		case uint64:
			checkArr = append(checkArr, fmt.Sprintf("%s=%d", key, int64(val.(uint64))))
		}

	}
	sort.Strings(checkArr)
	checkStr := strings.Join(checkArr, "\n")
	secretKey := utils.Sha256Bytes(botToken)
	if utils.Hmac256(checkStr, secretKey) != checkHash {
		return false
	}
	var authDate int64
	switch telegram["auth_date"].(type) {
	case float64:
		authDate = int64(telegram["auth_date"].(float64))
	case int:
		authDate = int64(telegram["auth_date"].(int))
	case int64:
		authDate = telegram["auth_date"].(int64)
	case uint64:
		authDate = int64(telegram["auth_date"].(uint64))
	}
	if time.Now().Unix()-authDate > 86400 {
		return false
	}
	return true
}
