package common

type Config struct {
	AppName              string        `default:"tokenmama"`
	BaseUrl              string        `default:"https://tokenmama.io"`
	CDNUrl               string        `default:"https://adxcdn.tokenmama.io/"`
	CreativeCDN          string        `default:"https://adcdn.tokenmama.io"`
	AdUrl                string        `default:"https://adx.tokenmama.io/t/"`
	AdImpUrl             string        `default:"https://adx.tokenmama.io/i/"`
	PublisherDomain      string        `default:"https://media.tokenmama.io"`
	AdvertiserDomain     string        `default:"https://adx.tokenmama.io"`
	CookieDomain         string        `default:"*.tokenmama.io"`
	Port                 int           `default:"8005"`
	Geth                 string        `default:"https://mainnet.infura.io/NlT37dDxuLT2tlZNw3It"`
	UI                   string        `required:"true"`
	Template             string        `required:"true"`
	LogPath              string        `required:"true"`
	TokenSalt            string        `required:"true"`
	LinkSalt             string        `required:"true"`
	SentryDSN            string        `required:"true"`
	Airdrop              AirdropConfig `required:"true"`
	AuctionRate          float64       `required:"true"`
	ClickhouseDSN        string        `required:"true"`
	MySQL                MySQLConfig   `required:"true"`
	Redis                RedisConfig   `required:"true"`
	SQS                  SQSConfig     `required:"true"`
	S3                   S3Config      `required:"true"`
	Mail                 MailConfig    `required:"true"`
	SlackToken           string        `required:"true"`
	SlackAdminChannelID  string        `required:"true"`
	TwilioToken          string        `required:"true"`
	TelegramBotToken     string        `required:"true"`
	TelegramBotName      string        `required:"true"`
	GeoIP                string        `required:"true"`
	AdJSVer              string        `required:"true"`
	Debug                bool          `default:"false"`
	EnableWeb            bool          `default:"false"`
	EnableTelegramBot    bool          `default:"false"`
	EnableGC             bool          `default:"false"`
	EnableDealer         bool          `default:"false"`
	EnableDepositChecker bool          `default:"false"`
	EnableAdServer       bool          `default:"false"`
	Zz253                Sms           `required:"true"`
}

type Sms struct {
	Account  string `default:"N5692616"`
	Password string `default:"4mNud0qth"`
}

type AirdropConfig struct {
	CommissionFee          uint64 `default:"4"`
	GasPrice               int64  `default:"5"`
	GasLimit               int64  `default:"210000"`
	DealerContractGasPrice int64  `default:"5"`
	DealerContractGasLimit uint64 `default:"210000"`
}

type SQSConfig struct {
	Region       string `default:"ap-northeast-1"`
	EmailQueue   string `default:"email"`
	AdClickQueue string `default:"adclick"`
	AdImpQueue   string `default:"adimp"`
	AccountId    string `required:"true"`
	AK           string `required:"true"`
	Secret       string `required:"true"`
	Token        string `default:""`
}

type S3Config struct {
	Region       string `default:"ap-northeast-1"`
	AK           string `required:"true"`
	Secret       string `required:"true"`
	AdBucket     string `required:"true"`
	CreativePath string `requird:"true"`
	Token        string `default:""`
}

type MySQLConfig struct {
	Host   string `required:"true"`
	User   string `required:"true"`
	Passwd string `required:"true"`
	DB     string `default:"tokenme"`
}

type RedisConfig struct {
	Master string `required:"true"`
	Slave  string
}

type MailConfig struct {
	Server string `default:"localhost"`
	Port   int    `default:"25"`
	User   string
	Passwd string
}
