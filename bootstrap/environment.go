package bootstrap

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Environment struct {
	DataBase   DataBase   `mapstructure:",squash" json:"database"`
	Redis      Redis      `mapstructure:",squash" json:"redis"`
	HttpServer HttpServer `mapstructure:",squash" json:"http_server"`
	JWT        JWT        `mapstructure:",squash" json:"jwt"`
	Storage    Storage    `mapstructure:",squash" json:"storage"`
	Limiter    Limiter    `mapstructure:",squash" json:"limiter"`
	OTP        OTP        `mapstructure:",squash" json:"otp"`
}

type HttpServer struct {
	Domain string `mapstructure:"HTTP_DOMAIN" json:"domain"`
	Port   string `mapstructure:"HTTP_PORT"   json:"port"`
}

type DataBase struct {
	Host     string `mapstructure:"PDB_HOST"     json:"host"`
	Port     string `mapstructure:"PDB_PORT"     json:"port"`
	User     string `mapstructure:"PDB_USER"     json:"user"`
	Password string `mapstructure:"PDB_PASSWORD" json:"password"`
	Name     string `mapstructure:"PDB_DB_NAME"  json:"name"`
	SSLMode  string `mapstructure:"PDB_SSL_MODE" json:"ssl_mode"`
}

type Redis struct {
	Host      string `mapstructure:"RDB_HOST"      json:"host"`
	Port      string `mapstructure:"RDB_PORT"      json:"port"`
	Password  string `mapstructure:"RDB_PASSWORD"  json:"password"`
	RDBNumber string `mapstructure:"RDB_DB_NUMBER" json:"db"`
}

type JWT struct {
	AccessPrivateKey    string `mapstructure:"JWT_ACCESS_PRIVATE"        json:"access_private_key"`
	AccessExpireMinutes int    `mapstructure:"JWT_ACCESS_EXPIRE_MINUTES" json:"access_expire_minutes"`

	RefreshPrivateKey string `mapstructure:"JWT_REFRESH_PRIVATE"        json:"refresh_private_key"`
	RefreshExpireDays int    `mapstructure:"JWT_REFRESH_EXPIRE_DAYS"   json:"refresh_expire_days"`
}

type Storage struct {
	AccessKey string     `mapstructure:"S3_ACCESS_KEY" json:"access_key"`
	SecretKey string     `mapstructure:"S3_SECRET_KEY" json:"secret_key"`
	Region    string     `mapstructure:"S3_REGION"     json:"region"`
	Endpoint  string     `mapstructure:"S3_ENDPOINT"   json:"endpoint"`
	Buckets   BucketName `mapstructure:",squash"       json:"buckets"`
}

type BucketName struct {
	MoviePoster  string `mapstructure:"S3_MOVIE_POSTER_BUCKET"  json:"movie_poster_bucket"`
	MovieTrailer string `mapstructure:"S3_MOVIE_TRAILER_BUCKET" json:"movie_trailer_bucket"`
	ActorPhoto   string `mapstructure:"S3_ACTOR_PHOTO_BUCKET"   json:"actor_photo_bucket"`
	CinemaPhoto  string `mapstructure:"S3_CINEMA_PHOTO_BUCKET"  json:"cinema_photo_bucket"`
}

type Limiter struct {
	HighRate   int `mapstructure:"LIMITER_HIGH_RATE"   json:"high_rate"`
	MediumRate int `mapstructure:"LIMITER_MEDIUM_RATE" json:"medium_rate"`
	SmallRate  int `mapstructure:"LIMITER_SMALL_RATE"  json:"small_rate"`
	TinyRate   int `mapstructure:"LIMITER_TINY_RATE"   json:"tiny_rate"`
}

type OTP struct {
	Digit         int64 `mapstructure:"OTP_DIGIT_NUMBER"   json:"digit_number"`
	ExpireMinutes int   `mapstructure:"OTP_EXPIRE_MINUTES" json:"expire_minutes"`
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

func NewEnvironment() *Environment {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("../..")
	v.SetConfigType("env")

	v.SetConfigName(".env")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
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
		return fmt.Sprintf("Environment{DataBase:%+v, Redis:%+v, HttpServer:%+v, JWT:%+v, Storage:%+v, Limiter:%+v, OTP:%+v}",
			e.DataBase, e.Redis, e.HttpServer, e.JWT, e.Storage, e.Limiter, e.OTP)
	}

	return string(data)
}
