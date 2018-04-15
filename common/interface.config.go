package common

type Config struct {
	AppName             string      `default:"tokenmama"`
	BaseUrl             string      `default:"https://tokenmama.io"`
	CDNUrl              string      `default:"https://static.tianxi100.com/"`
	PublisherDomain     string      `default:"https://media.tokenmama.io"`
	AdvertiserDomain    string      `default:"https://adx.tokenmama.io"`
	Port                int         `default:"8005"`
	UI                  string      `default:"./ui/dist"`
	LogPath             string      `default:"/tmp/tokenmama-adx"`
	Debug               bool        `default:"false"`
	MySQL               MySQLConfig `required:"true"`
	Redis               RedisConfig
	Geth                string    `default:"https://mainnet.infura.io/NlT37dDxuLT2tlZNw3It"`
	SlackToken          string    `required:"true"`
	SlackAdminChannelID string    `default:"G9Y7METUG"`
	TelegramBotToken    string    `required:"true"`
	TelegramBotName     string    `required:"true"`
	GeoIP               string    `required:"true"`
	TokenSalt           string    `default:"20eefe8d82ba3ca8a417e14a48d24632bc35bbd7"`
	LinkSalt            string    `default:"20eefe8d82ba3ca8a417e14a48d24632bc35bbd7"`
	SentryDSN           string    `default:"https://6a576b1028974e93a5e2d29071c0e896:fc012faed5a94a3683f435714e1dc2e1@sentry.io/1163662"`
	SQS                 SQSConfig `required:"true"`
	EnableWeb           bool
	EnableGC            bool
	Mail                MailConfig `required:"true"`
}

type SQSConfig struct {
	AccountId   string `required:"true"`
	AK          string `required:"true"`
	Secret      string `required:"true"`
	Token       string
	EmailQueue  string `default:"email"`
	EmailRegion string `default:"ap-northeast-1"`
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
