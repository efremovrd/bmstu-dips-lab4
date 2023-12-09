package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	// Cors     CorsConfig
}

type ServerConfig struct {
	AppVersion   string
	Port         string
	JwtSecretKey string
	CtxUserKey   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
}

// type CorsConfig struct {
// 	AllowOrigins []string
// 	AllowMethods []string
// 	AllowHeaders []string
// }

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
