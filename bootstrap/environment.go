package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Environment struct {
	DebugModeOptions    DebugModeOptions    `mapstructure:",squash" json:"debug_mode_options"`
	DataBase            DataBase            `mapstructure:",squash" json:"database"`
	Redis               Redis               `mapstructure:",squash" json:"redis"`
	HttpServer          HttpServer          `mapstructure:",squash" json:"http_server"`
	JWT                 JWT                 `mapstructure:",squash" json:"jwt"`
	Storage             Storage             `mapstructure:",squash" json:"storage"`
	Limiter             Limiter             `mapstructure:",squash" json:"limiter"`
	OTP                 OTP                 `mapstructure:",squash" json:"otp"`
	PEM                 PEM                 `mapstructure:",squash" json:"pem"`
	UserPermissionCache UserPermissionCache `mapstructure:",squash" json:"user_permission_cache"`
	Logger              LoggerConfig        `mapstructure:",squash" json:"logger"`
	Email               EmailConfig         `mapstructure:",squash" json:"email"`
	Captcha             CaptchaConfig       `mapstructure:",squash" json:"captcha"`
	Elasticsearch       ElasticsearchConfig `mapstructure:",squash" json:"elasticsearch"`
	RabbitMQ            RabbitMQConfig      `mapstructure:",squash" json:"rabbitmq"`
	Reservation         ReservationConfig   `mapstructure:",squash" json:"reservation"`
	AutoMigrate         bool                `mapstructure:"AUTO_MIGRATE" json:"auto_migrate"`
	CORSAllowedOrigins  string              `mapstructure:"CORS_ALLOWED_ORIGINS" json:"cors_allowed_origins"`
	SuperAdminPhone     string              `mapstructure:"SUPER_ADMIN_PHONE" json:"-"`
	SuperAdminEmail     string              `mapstructure:"SUPER_ADMIN_EMAIL" json:"-"`
	SuperAdminPassword  string              `mapstructure:"SUPER_ADMIN_PASSWORD" json:"-"`
	SuperAdminFirstName string              `mapstructure:"SUPER_ADMIN_FIRST_NAME" json:"-"`
	SuperAdminLastName  string              `mapstructure:"SUPER_ADMIN_LAST_NAME" json:"-"`
	SuperAdminWallet    int64               `mapstructure:"SUPER_ADMIN_WALLET" json:"-"`
}

type DebugModeOptions struct {
	Enable bool `mapstructure:"DEBUG_MODE_ENABLE" json:"enable_debug_mode"`
}

type HttpServer struct {
	Domain string `mapstructure:"HTTP_DOMAIN" json:"domain"`
	Port   string `mapstructure:"HTTP_PORT"   json:"port"`
}

type DataBase struct {
	Host     string `mapstructure:"PDB_HOST"     json:"host"`
	Port     string `mapstructure:"PDB_PORT"     json:"port"`
	User     string `mapstructure:"PDB_USER"     json:"user"`
	Password string `mapstructure:"PDB_PASSWORD" json:"-"`
	Name     string `mapstructure:"PDB_DB_NAME"  json:"name"`
	SSLMode  string `mapstructure:"PDB_SSL_MODE" json:"ssl_mode"`
}

type Redis struct {
	Host      string `mapstructure:"RDB_HOST"      json:"host"`
	Port      string `mapstructure:"RDB_PORT"      json:"port"`
	Password  string `mapstructure:"RDB_PASSWORD"  json:"-"`
	RDBNumber string `mapstructure:"RDB_DB_NUMBER" json:"db"`
}

type PEM struct {
	PrivateAccessPath  string `json:"private_access_path" mapstructure:"PRIVATE_ACCESS_PEM_PATH"`
	PrivateRefreshPath string `json:"private_refresh_path" mapstructure:"PRIVATE_REFRESH_PEM_PATH"`
	PublicAccessPath   string `json:"public_access_path" mapstructure:"PUBLIC_ACCESS_PEM_PATH"`
	PublicRefreshPath  string `json:"public_refresh_path" mapstructure:"PUBLIC_REFRESH_PEM_PATH"`
}

type JWT struct {
	AccessExpireMinutes int `mapstructure:"JWT_ACCESS_EXPIRE_MINUTES" json:"access_expire_minutes"`
	RefreshExpireDays   int `mapstructure:"JWT_REFRESH_EXPIRE_DAYS"   json:"refresh_expire_days"`
}

type Storage struct {
	AccessKey string     `mapstructure:"S3_ACCESS_KEY" json:"-"`
	SecretKey string     `mapstructure:"S3_SECRET_KEY" json:"-"`
	Region    string     `mapstructure:"S3_REGION"     json:"region"`
	Endpoint  string     `mapstructure:"S3_ENDPOINT"   json:"endpoint"`
	Buckets   BucketName `mapstructure:",squash"       json:"buckets"`
}

type BucketName struct {
	MoviePoster  string `mapstructure:"S3_MOVIE_POSTER_BUCKET"  json:"movie_poster_bucket"`
	MovieTrailer string `mapstructure:"S3_MOVIE_TRAILER_BUCKET" json:"movie_trailer_bucket"`
	ActorPhoto   string `mapstructure:"S3_ACTOR_PHOTO_BUCKET"   json:"actor_photo_bucket"`
	CinemaPhoto  string `mapstructure:"S3_CINEMA_PHOTO_BUCKET"  json:"cinema_photo_bucket"`
	TicketPDF    string `mapstructure:"S3_TICKET_PDF_BUCKET"    json:"ticket_pdf_bucket"`
}

type Limiter struct {
	HighRate   int `mapstructure:"LIMITER_HIGH_RATE"   json:"high_rate"`
	MediumRate int `mapstructure:"LIMITER_MEDIUM_RATE" json:"medium_rate"`
	SmallRate  int `mapstructure:"LIMITER_SMALL_RATE"  json:"small_rate"`
	TinyRate   int `mapstructure:"LIMITER_TINY_RATE"   json:"tiny_rate"`
}

type OTP struct {
	Digit               int64 `mapstructure:"OTP_DIGIT_NUMBER"   json:"digit_number"`
	ExpireMinutes       int   `mapstructure:"OTP_EXPIRE_MINUTES" json:"expire_minutes"`
	VerifyExpireMinutes int   `mapstructure:"OTP_VERIFY_EXPIRE_MINUTES" json:"verify_expire_minutes"`
}

type LoggerConfig struct {
	FilePath string `json:"file_path" mapstructure:"LOGGER_FILE_PATH"`
	Encoding string `json:"encoding" mapstructure:"LOGGER_ENCODING"`
	Level    string `json:"level" mapstructure:"LOGGER_LEVEL"`
	Logger   string `json:"logger" mapstructure:"LOGGER_NAME"`
}

type EmailConfig struct {
	CinemaTicketEmail            string `mapstructure:"CINEMA_TICKET_EMAIL"              json:"cinema_ticket_email"`
	CinemaTicketEmailAppPassword string `mapstructure:"CINEMA_TICKET_EMAIL_APP_PASSWORD" json:"-"`
	SMTPHost                     string `mapstructure:"SMTP_HOST"                        json:"smtp_host"`
	SMTPPort                     string `mapstructure:"SMTP_HOST_PORT"                   json:"smtp_host_port"`
}

type CaptchaConfig struct {
	Type string `mapstructure:"CAPTCHA_TYPE" json:"type"`
	TTL  int    `mapstructure:"CAPTCHA_TTL_IN_SECONDS" json:"ttl"`
}

type UserPermissionCache struct {
	UserPermissionExpiredDuration int `mapstructure:"USER_PERMISSION_CACHE_EXPIRE_MINUTES" json:"expire_duration"`
}

type ElasticsearchConfig struct {
	URL string `mapstructure:"ELASTICSEARCH_URL" json:"url"`
}

type RabbitMQConfig struct {
	URL string `mapstructure:"RABBITMQ_URL" json:"-"`
}

type ReservationConfig struct {
	HoldSeconds         int `mapstructure:"RESERVATION_HOLD_SECONDS" json:"hold_seconds"`
	RefundCutoffMinutes int `mapstructure:"REFUND_CUTOFF_MINUTES" json:"refund_cutoff_minutes"`
}

func (r *ReservationConfig) RefundCutoff() time.Duration {
	minutes := r.RefundCutoffMinutes
	if minutes <= 0 {
		minutes = 120
	}
	return time.Duration(minutes) * time.Minute
}

func (r *ReservationConfig) HoldDuration() time.Duration {
	seconds := r.HoldSeconds
	if seconds <= 0 {
		seconds = 600
	}
	return time.Duration(seconds) * time.Second
}

func (c *JWT) AccessDuration() time.Duration {
	return time.Duration(c.AccessExpireMinutes) * time.Minute
}

func (c *JWT) RefreshDuration() time.Duration {
	return time.Duration(c.RefreshExpireDays) * 24 * time.Hour
}

func (o *OTP) ExpirationDuration() time.Duration {
	return time.Duration(o.ExpireMinutes) * time.Minute
}
func (o *OTP) ExpirationDurationInMinutes() int {
	return int((time.Duration(o.ExpireMinutes) * time.Minute)) / int(time.Minute)
}

func (o *UserPermissionCache) ExpirationDuration() time.Duration {
	return time.Duration(o.UserPermissionExpiredDuration) * time.Minute
}

func (o *OTP) VerifyExpirationDuration() time.Duration {
	return time.Duration(o.VerifyExpireMinutes) * time.Minute
}

func (cc *CaptchaConfig) TTLInSeconds() time.Duration {
	return time.Duration(cc.TTL) * time.Second
}

func NewEnvironment() *Environment {
	v := viper.New()
	v.SetDefault("RESERVATION_HOLD_SECONDS", 600)
	v.SetDefault("REFUND_CUTOFF_MINUTES", 120)
	v.SetDefault("AUTO_MIGRATE", true)
	for _, key := range []string{
		"DEBUG_MODE_ENABLE", "HTTP_DOMAIN", "HTTP_PORT",
		"PDB_HOST", "PDB_PORT", "PDB_USER", "PDB_PASSWORD", "PDB_DB_NAME", "PDB_SSL_MODE",
		"RDB_HOST", "RDB_PORT", "RDB_PASSWORD", "RDB_DB_NUMBER",
		"JWT_ACCESS_EXPIRE_MINUTES", "JWT_REFRESH_EXPIRE_DAYS",
		"S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_REGION", "S3_ENDPOINT",
		"S3_MOVIE_POSTER_BUCKET", "S3_MOVIE_TRAILER_BUCKET", "S3_ACTOR_PHOTO_BUCKET", "S3_CINEMA_PHOTO_BUCKET", "S3_TICKET_PDF_BUCKET",
		"LIMITER_HIGH_RATE", "LIMITER_MEDIUM_RATE", "LIMITER_SMALL_RATE", "LIMITER_TINY_RATE",
		"OTP_DIGIT_NUMBER", "OTP_EXPIRE_MINUTES", "OTP_VERIFY_EXPIRE_MINUTES",
		"PRIVATE_ACCESS_PEM_PATH", "PRIVATE_REFRESH_PEM_PATH", "PUBLIC_ACCESS_PEM_PATH", "PUBLIC_REFRESH_PEM_PATH",
		"USER_PERMISSION_CACHE_EXPIRE_MINUTES", "LOGGER_FILE_PATH", "LOGGER_ENCODING", "LOGGER_LEVEL", "LOGGER_NAME",
		"CINEMA_TICKET_EMAIL", "CINEMA_TICKET_EMAIL_APP_PASSWORD", "SMTP_HOST", "SMTP_HOST_PORT",
		"CAPTCHA_TTL_IN_SECONDS", "CAPTCHA_TYPE", "ELASTICSEARCH_URL", "RABBITMQ_URL",
		"RESERVATION_HOLD_SECONDS", "REFUND_CUTOFF_MINUTES", "CORS_ALLOWED_ORIGINS", "AUTO_MIGRATE",
		"SUPER_ADMIN_PHONE", "SUPER_ADMIN_EMAIL", "SUPER_ADMIN_PASSWORD",
		"SUPER_ADMIN_FIRST_NAME", "SUPER_ADMIN_LAST_NAME", "SUPER_ADMIN_WALLET",
	} {
		_ = v.BindEnv(key)
	}

	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("../..")
	v.SetConfigType("env")

	v.SetConfigName(".env")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			panic(fmt.Errorf("failed to read config file: %w", err))
		}
	}

	var env Environment
	if err := v.Unmarshal(&env); err != nil {
		panic(fmt.Errorf("failed to unmarshal environment config: %w", err))
	}

	return &env
}

func (e *Environment) String() string {
	if e == nil {
		return "<nil>"
	}

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "Environment{redacted: true}"
	}

	return string(data)
}
