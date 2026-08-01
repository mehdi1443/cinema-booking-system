package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Environment struct {
	DataBase   DataBase   `mapstructure:",squash" json:"database"`
	Redis      Redis      `mapstructure:",squash" json:"redis"`
	HttpServer HttpServer `mapstructure:",squash" json:"http_server"`
	JWT        JWT        `mapstructure:",squash" json:"jwt"`
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
	Name     string `mapstructure:"PDB_NAME"     json:"name"`
	SSLMode  string `mapstructure:"PDB_SSL_MODE" json:"ssl_mode"`
}

type Redis struct {
	Host      string `mapstructure:"RDB_HOST"     json:"host"`
	Port      string `mapstructure:"RDB_PORT"     json:"port"`
	Password  string `mapstructure:"RDB_PASSWORD" json:"password"`
	RDBNumber string `mapstructure:"RDB_NUMBER"   json:"rdb_number"`
}

type JWT struct {
	AccessSecret        string `mapstructure:"JWT_ACCESS_SECRET"         json:"access_secret"`
	RefreshSecret       string `mapstructure:"JWT_REFRESH_SECRET"        json:"refresh_secret"`
	AccessExpireMinutes int    `mapstructure:"JWT_ACCESS_EXPIRE_MINUTES" json:"access_expire_minutes"`
	RefreshExpireDays   int    `mapstructure:"JWT_REFRESH_EXPIRE_DAYS"   json:"refresh_expire_days"`
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

// TODO : Declare String() Func for Debug
func (e *Environment) String() string {
	if e == nil {
		return "<nil>"
	}

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("Environment{DataBase:%+v, Redis:%+v, HttpServer:%+v, JWT:%+v}",
			e.DataBase, e.Redis, e.HttpServer, e.JWT)
	}

	return string(data)
}
