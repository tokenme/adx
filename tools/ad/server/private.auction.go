package server

import (
	"errors"
	//"github.com/davecgh/go-spew/spew"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	"github.com/tokenme/adx/utils"
	"sync"
)

type PrivateAuctionPool struct {
	ads     map[uint64][][]common.Ad
	adzones map[uint64]common.Adzone
	sync.RWMutex
}

func NewPrivateAuctionPool() *PrivateAuctionPool {
	return &PrivateAuctionPool{
		ads:     make(map[uint64][][]common.Ad),
		adzones: make(map[uint64]common.Adzone),
	}
}

func (this *PrivateAuctionPool) Reload(service *common.Service, config common.Config) {
	db := service.Db
	query := `SELECT
    pac.id ,
    pac.auction_id ,
    pac.adzone_id ,
    pac.media_id ,
    pac.user_id AS advertiser_id ,
    a.user_id AS publisher_id ,
    a.size_id ,
    pac.landing_page ,
    pac.img ,
    a.placeholder_url, 
    a.placeholder_img,
    a.user_id
FROM
    adx.private_auction_creatives AS pac
INNER JOIN adx.private_auctions AS pa ON ( pa.id = pac.auction_id )
INNER JOIN adx.adzones AS a ON ( a.id = pac.adzone_id )
WHERE
    pa.online_status = 1
AND pa.audit_status = 1
AND pa.start_on <= DATE( NOW())
AND pa.end_on >= DATE( NOW())`
	rows, _, err := db.Query(query)
	if err != nil {
		log.Error(err.Error())
		return
	}
	adsMap := make(map[uint64]map[uint64][]common.Ad)
	adzones := make(map[uint64]common.Adzone)
	for _, row := range rows {
		ad := common.Ad{
			Id:           row.Uint64(0),
			AuctionId:    row.Uint64(1),
			AdzoneId:     row.Uint64(2),
			MediaId:      row.Uint64(3),
			AdvertiserId: row.Uint64(4),
			PublisherId:  row.Uint64(5),
			SizeId:       row.Uint64(6),
			Url:          row.Str(7),
			Img:          row.Str(8),
		}
		ad.ImgUrl = ad.GetImgUrl(config)
		if _, found := adsMap[ad.AdzoneId]; !found {
			adsMap[ad.AdzoneId] = make(map[uint64][]common.Ad)
		}
		adsMap[ad.AdzoneId][ad.AuctionId] = append(adsMap[ad.AdzoneId][ad.AuctionId], ad)
		placeholderUrl := row.Str(9)
		placeholderImg := row.Str(10)
		if placeholderImg == "" || placeholderUrl == "" {
			continue
		}
		if _, found := adzones[ad.AdzoneId]; found {
			continue
		}

		adzones[ad.AdzoneId] = common.Adzone{
			Id: ad.AdzoneId,
			Media: common.Media{
				Id: ad.MediaId,
			},
			UserId: row.Uint64(11),
			Size: common.Size{
				Id: uint(ad.SizeId),
			},
			Placeholder: &common.PrivateAuctionCreative{
				Url: placeholderUrl,
				Img: placeholderImg,
			},
		}
	}
	ads := make(map[uint64][][]common.Ad)
	for adzoneId, auctionMap := range adsMap {
		for _, creatives := range auctionMap {
			ads[adzoneId] = append(ads[adzoneId], creatives)
		}
	}
	this.Lock()
	this.ads = ads
	this.adzones = adzones
	this.Unlock()
}

func (this *PrivateAuctionPool) Pop(adzoneId uint64) (ad common.Ad, err error) {
	this.RLock()
	ads := this.ads[adzoneId]
	adzone, adzoneFound := this.adzones[adzoneId]
	this.RUnlock()
	if ads == nil {
		if !adzoneFound {
			return ad, errors.New("not found")
		}
		return common.Ad{
			AdvertiserId:  adzone.UserId,
			PublisherId:   adzone.UserId,
			AdzoneId:      adzone.Id,
			MediaId:       adzone.Media.Id,
			SizeId:        uint64(adzone.Size.Id),
			Url:           adzone.Placeholder.Url,
			Img:           adzone.Placeholder.Img,
			IsPlaceholder: 1,
		}, nil
	}
	var idx int
	totalAds := len(ads)
	if totalAds > 1 {
		idx = utils.RangeRandInt(0, totalAds-1)
	}
	auctions := ads[idx]
	totalAuctions := len(auctions)
	if totalAuctions > 1 {
		idx := utils.RangeRandInt(0, totalAuctions-1)
		ad = auctions[idx]
	} else {
		ad = auctions[0]
	}
	ad.PK, _ = ad.GetPK()
	return ad, nil
}
