package common

import (
	"fmt"
	"github.com/tokenme/adx/utils"
	"github.com/tokenme/adx/utils/binary"
)

type Ad struct {
	PK            string `json:"pk"`
	LogTime       int64  `json:"log_time"`
	Id            uint64 `json:"id"`
	AuctionId     uint64 `json:"auction_id,omitempty"`
	AdzoneId      uint64 `json:"adzone_id,omitempty"`
	MediaId       uint64 `json:"media_id,omitempty"`
	SizeId        uint64 `json:"size_id,omitempty"`
	AdvertiserId  uint64 `json:"advertiser_id,omitempty"`
	PublisherId   uint64 `jsonn:"publisher_id,omitempty"`
	Url           string `json:"url,omitempty"`
	Img           string `json:"img,omitempty"`
	ImgUrl        string `json:"img_url,omitempty"`
	IsPlaceholder uint   `json:"is_placeholder,omitempty"`
	Env           AdEnv  `json:"env,omitempty"`
}

type AdEnv struct {
	Cookie         string `json:"cookie,omitempty"`
	URL            string `json:"url,omitempty"`
	Referrer       string `json:"referrer,omitempty"`
	ScreenSize     string `json:"screen_size,omitempty"`
	AdSize         string `json:"ad_size,omitempty"`
	OsName         string `json:"os_name,omitempty"`
	OsVersion      string `json:"os_version,omitempty"`
	BrowserName    string `json:"browser_name,omitempty"`
	BrowserVersion string `json:"browser_version,omitempty"`
	BrowserType    uint   `json:"browser_type,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
	IP             int64  `json:"ip,omitempty"`
	CountryId      uint   `json:"country_id,omitempty"`
	CountryName    string `json:"country_name,omitempty"`
}

func (this Ad) GetPK() (string, error) {
	return utils.Salt()
}

func (this Ad) GetImgUrl(config Config) string {
	return fmt.Sprintf("%s/%s/%s", config.CreativeCDN, config.S3.CreativePath, this.Img)
}

func (this Ad) GetLink(config Config) string {
	encoded, _ := EncodeAd([]byte(config.LinkSalt), this)
	return fmt.Sprintf("%s%s", config.AdUrl, encoded)
}

func (this Ad) GetImpUrl(config Config) string {
	encoded, _ := EncodeAd([]byte(config.LinkSalt), this)
	return fmt.Sprintf("%s%s", config.AdImpUrl, encoded)
}

func EncodeAd(key []byte, ad Ad) (string, error) {
	enc := binary.NewEncoder()
	enc.Encode(ad)
	return utils.AESEncryptBytes(key, enc.Buffer())
}

func DecodeAd(key []byte, cryptoText string) (ad Ad, err error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return ad, err
	}
	dec := binary.NewDecoder()
	dec.SetBuffer(data)
	dec.Decode(&ad)
	return ad, nil
}
