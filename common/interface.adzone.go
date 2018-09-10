package common

import (
	"fmt"
	"github.com/tokenme/adx/utils"
	"strings"
	"time"
)

type Settlement = uint

const (
	CPM  Settlement = 1
	CPT  Settlement = 2
	CPMT Settlement = 3
)

type Adzone struct {
	Id                 uint64                  `json:"id"`
	UserId             uint64                  `json:"user_id,omitempty"`
	Media              Media                   `json:"media"`
	Size               Size                    `json:"size"`
	Url                string                  `json:"url"`
	MinCPT             float64                 `json:"min_cpt,omitempty"`
	MinCPM             float64                 `json:"min_cpm,omitempty"`
	SuggestCPT         float64                 `json:"suggest_cpt,omitempty"`
	Settlement         Settlement              `json:"settlement"`
	Desc               string                  `json:"desc"`
	Rolling            uint                    `json:"rolling"`
	OnlineStatus       uint                    `json:"online_status"`
	Placeholder        *PrivateAuctionCreative `json:"placeholder"`
	UnverifiedAuctions uint                    `json:"unverified_auctions,omitempty"`
	UnavailableDays    []time.Time             `json:"unavailable_days,omitempty"`
	EmbedCode          string                  `json:"embed,omitempty"`
	InsertedAt         time.Time               `json:"inserted_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
	Advantage          string                  `json:"advantage"`
	Location           string                  `json:"location"`
	Traffic            string                  `json:"traffic"`
}

func EncodeAdzoneId(key []byte, id uint64) (string, error) {
	buf := utils.Uint64ToByte(id)
	return utils.AESEncryptBytes(key, buf)
}

func DecodeAdzoneId(key []byte, cryptoText string) (uint64, error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return 0, err
	}
	return utils.ByteToUint64(data), nil
}

func (this Adzone) GetEmbedCode(config Config) string {
	embed := `<script type="text/javascript">
  var tmm_id='%s', tmm_width=%d, tmm_height=%d;
</script>
<script src="https://adx.tokenmama.io/tmm.js" type="text/javascript"></script>`
	encodedId, _ := EncodeAdzoneId([]byte(config.LinkSalt), this.Id)
	return fmt.Sprintf(embed, encodedId, this.Size.Width, this.Size.Height)
}

func (this Adzone) GetUnavailableDays(service *Service) (days []time.Time, err error) {
	db := service.Db
	query := `SELECT
	adzone_id ,
	COUNT( aad.auction_id ) AS auctions ,
	a.rolling AS rolling ,
	aad.record_on
FROM
	adx.adzone_auction_days AS aad
INNER JOIN adx.adzones AS a ON ( a.id = aad.adzone_id )
WHERE aad.adzone_id=%d AND record_on >= DATE(NOW())
GROUP BY
	adzone_id ,
	record_on
HAVING
	auctions >= rolling`
	rows, _, err := db.Query(query, this.Id)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		days = append(days, row.ForceLocaltime(3))
	}
	return days, nil
}

func AdzonesUnavailableDays(service *Service, adzoneIds []string) (unavaliableDays map[uint64][]time.Time, err error) {
	if len(adzoneIds) == 0 {
		return nil, nil
	}
	db := service.Db
	query := `SELECT
	adzone_id ,
	COUNT( aad.auction_id ) AS auctions ,
	a.rolling AS rolling ,
	aad.record_on
FROM
	adx.adzone_auction_days AS aad
INNER JOIN adx.adzones AS a ON ( a.id = aad.adzone_id )
WHERE aad.adzone_id IN (%s) AND record_on >= DATE(NOW())
GROUP BY
	adzone_id ,
	record_on
HAVING
	auctions >= rolling`
	unavaliableDays = make(map[uint64][]time.Time)
	rows, _, err := db.Query(query, strings.Join(adzoneIds, ","))
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		adzoneId := row.Uint64(0)
		unavaliableDays[adzoneId] = append(unavaliableDays[adzoneId], row.ForceLocaltime(3))
	}
	return unavaliableDays, nil
}

func UnavailableAdzones(service *Service, dateRange []time.Time) (adzoneIds []uint64, err error) {
	if len(dateRange) != 2 || dateRange[1].Before(dateRange[0]) {
		dateRange[0] = utils.TimeToDate(time.Now())
		dateRange[1] = dateRange[0].AddDate(0, 2, 0)
	}
	days := int(dateRange[1].Sub(dateRange[0]).Hours())/24 + 1
	db := service.Db
	query := `SELECT
	adzone_id ,
	COUNT(*) AS num
FROM
	( SELECT
		adzone_id ,
		COUNT( aad.auction_id ) AS auctions ,
		a.rolling AS rolling ,
		aad.record_on AS record_on
	FROM
		adx.adzone_auction_days AS aad
	INNER JOIN adx.adzones AS a ON ( a.id = aad.adzone_id )
	WHERE
		record_on BETWEEN '%s' AND '%s'
	GROUP BY
		adzone_id ,
		record_on
	HAVING
		auctions >= rolling ) AS tmp
GROUP BY
	tmp.adzone_id
HAVING num >= %d`
	rows, _, err := db.Query(query, dateRange[0].Format("2006-01-02"), dateRange[1].Format("2006-01-02"), days)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		adzoneIds = append(adzoneIds, row.Uint64(0))
	}
	return adzoneIds, nil
}
