package singleton

import (
	"log"

	"github.com/ducthangng/GeoFleet/service/copier"
	"github.com/spf13/viper"
)

type GatewayGlobalConfig struct {
	// Server Connection
	Host           string `mapstructure:"HOST"`
	Port           string `mapstructure:"PORT"`
	RequestTimeout int    `mapstructure:"REQUEST_TIMEOUT"`
	RateLimit      int    `mapstructure:"RATE_LIMIT"`
	Env            string `mapstructure:"ENV"`

	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisUser string `mapstructure:"REDIS_USER"`
	RedisPass string `mapstructure:"REDIS_PASS"`
	RedisDB   int    `mapstructure:"REDIS_DB"`

	// Auth & Security
	Domain       string   `mapstructure:"DOMAIN"`
	JwtSecretKey string   `mapstructure:"JWT_SECRET_KEY"`
	AllowOrigins []string `mapstructure:"ALLOW_ORIGINS"`

	// Consul: Service Discovery & Dynamic Configuration
	ConsulAddress   string `mapstructure:"CONSUL_ADDRESS"`
	ServiceName     string `mapstructure:"SERVICE_NAME"`
	ServiceCheckURL string `mapstructure:"SERVICE_CHECK_URL"`

	// Kafka: Event Streaming / Audit Logging
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic   string   `mapstructure:"KAFKA_TOPIC"`
	KafkaGroupId string   `mapstructure:"KAFKA_GROUP_ID"`

	// Observability & Environment
	Environment string `mapstructure:"ENV"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	// Database
	Type            string `mapstructure:"TYPE"`
	User            string `mapstructure:"DB_USER"`
	Password        string `mapstructure:"PASSWORD"`
	DatabaseHost    string `mapstructure:"DBHOST"`
	Name            string `mapstructure:"NAME"`
	DatabasePort    string `mapstructure:"DBPORT"`
	SSLMode         string `mapstructure:"SSL_MODE"`
	CACERTBASE64    string `mapstructure:"CACERT_BASE64"`
	MaxIdleConns    int    `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns    int    `mapstructure:"MAX_OPEN_CONNS"`
	MaxConnLifeTime int    `mapstructure:"MAX_CONN_LIFE_TIME"`
	MaxConnIdleTime int    `mapstructure:"MAX_CONN_IDLE_TIME"`
}

var GlobalConfig *GatewayGlobalConfig

func InitializeConfig() {
	if GlobalConfig != nil {
		return
	}

	viper.AddConfigPath("./")
	viper.SetConfigType("env")
	viper.SetConfigName("tracking.env")
	viper.AutomaticEnv()

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		panic(err)
	}

	log.Println("config done")
}

func GetGlobalConfig() GatewayGlobalConfig {
	if GlobalConfig == nil {
		InitializeConfig()
	}

	var tempt GatewayGlobalConfig
	copier.MustCopy(&tempt, GlobalConfig)

	return tempt
}
