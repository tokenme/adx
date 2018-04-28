package common

import (
	"fmt"
)

type PrivateAuctionCreative struct {
	Id        uint64 `json:"id"`
	AuctionId uint64 `json:"auction_id,omitempty"`
	Url       string `json:"url,omitempty"`
	Img       string `json:"img,omitempty"`
	ImgUrl    string `json:"img_url,omitempty"`
}

func (this PrivateAuctionCreative) GetImgUrl(config Config) string {
	return fmt.Sprintf("%s/%s/%s", config.CreativeCDN, config.S3.CreativePath, this.Img)
}
