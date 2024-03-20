package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	RedisAddr string
	Burst     int64
	Rate      float64
	Window    time.Duration
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("error config file: %w", err))
		viper.SetDefault("REDIS_ADDR", "localhost:6379")
		viper.SetDefault("BURST", 16)
		viper.SetDefault("RATE", 16)
		viper.SetDefault("WINDOW", "8s")
	}
	return &Config{
		RedisAddr: viper.GetString("REDIS_ADDR"),
		Burst:     viper.GetInt64("BURST"),
		Rate:      viper.GetFloat64("RATE"),
		Window:    viper.GetDuration("WINDOW"),
	}
}
