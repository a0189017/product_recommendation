package config

import (
	"log"
	"os"
	"sync"

	"github.com/jinzhu/configor"
)

type Server struct {
	Host            string   `env:"HOST" default:"0.0.0.0"`
	Port            string   `env:"PORT" default:"8080"`
	EnableRotateLog bool     `env:"ENABLE_ROTATE_LOG_FILE" yaml:"enable_rotate_log"`
	EnableCORS      bool     `env:"ENABLE_CORS" default:"true" yaml:"enable_cors"`
	AllowOrigins    []string `env:"ALLOW_ORIGINS" yaml:"allow_origins"`
	UserAgent       string
}

type Settings struct {
	JWTSecret            string `env:"JWT_SECRET" yaml:"jwt_secret"`
	AccessTokenExpireSec string `env:"ACCESS_TOKEN_EXPIRES_IN_SECONDS" default:"14400" yaml:"access_token_expire_sec"`
	DataEncryptionKey    string `env:"DATA_ENCRYPTION_KEY" yaml:"data_encryption_key"`
	Concurrency          int    `env:"CONCURRENCY" default:"100"`
	OtpLength            int    `env:"OTP_LENGTH" default:"6"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST" yaml:"host"`
	Port     string `env:"REDIS_PORT" yaml:"port"`
	Password string `env:"REDIS_PASSWORD" yaml:"password"`
	DB       int    `env:"REDIS_DB" yaml:"db"`
}

type Config struct {
	Settings Settings
	DB       struct {
		User          string `env:"DB_USER" default:"admin"`
		Password      string `env:"DB_PASSWORD"`
		Host          string `env:"DB_HOST" default:"localhost"`
		Port          string `env:"DB_PORT" default:"5432"`
		Name          string `env:"DB_NAME" default:"db"`
		SlowThreshold string `env:"DB_SLOW_THRESHOLD" default:"1000" yaml:"slow_threshold"`
		Log           struct {
			Colorful string `env:"DB_LOG_COLORFUL" default:"false"`
		}
	}
	Redis   Redis
	Server  Server
	LogFile struct {
		FileName   string `env:"LOG_FILE_NAME" default:"temp.log" yaml:"file_name"`
		MaxSize    int    `env:"LOG_MAX_SIZE" default:"100" yaml:"max_size"`
		MaxBackups int    `env:"LOG_MAX_BACKUPS" default:"30" yaml:"max_backups"`
		MaxAge     int    `env:"LOG_MAX_AGE" default:"30" yaml:"max_age"`
	} `yaml:"log_file"`
	EmailService struct {
		Host           string `env:"EMAIL_HOST"`
		Account        string `env:"EMAIL_ACCOUNT"`
		Password       string `env:"EMAIL_PASSWORD"`
		DisplayName    string `env:"EMAIL_DISPLAY_NAME" yaml:"display_name"`
		SubjectPrefix  string `env:"EMAIL_SUBJECT_PREFIX" yaml:"subject_prefix"`
		TlsSupport     bool   `env:"EMAIL_TLS_SUPPORT" yaml:"tls_support"`
		EncryptionMode string `env:"EMAIL_ENCRYPTION_MODE" yaml:"encryption_mode"`
		FromEmail      string `env:"EMAIL_FROM_EMAIL" yaml:"from_email"`
	} `yaml:"email_service"`
	CronJobSpec struct {
		SyncRecommendation string
	}
}

var once sync.Once
var config *Config

func GetConfig(vaspCode ...string) *Config {
	once.Do(func() {
		// load config from env or yaml file
		// if config file does not exist, will just load from env
		// TODO: should be options
		configFilePath := os.Getenv("CONFIG_FILE_PATH")
		if configFilePath == "" {
			configFilePath = "config.yml"
		}
		config = &Config{}
		err := configor.Load(config, configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		setCornJobConfig()
	})

	return config
}

func setCornJobConfig() {
	config.CronJobSpec.SyncRecommendation = "55 */9 * * * *" // every 10 minutes at 55 seconds
}
