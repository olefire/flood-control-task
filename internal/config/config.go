package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	RedisAddr string
	N         int
	K         int
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("error config file: %w", err))
		viper.SetDefault("REDIS_ADDR", "localhost:6379")
		viper.SetDefault("N", 4)
		viper.SetDefault("K", 16)
	}
	return &Config{
		RedisAddr: viper.GetString("REDIS_ADDR"),
		N:         viper.GetInt("N"),
		K:         viper.GetInt("K"),
	}
}
