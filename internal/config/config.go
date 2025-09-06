package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram   TelegramConfig   `mapstructure:"telegram"`
	Models     ModelsConfig     `mapstructure:"api"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Ratelimit  RatelimitConfig  `mapstructure:"ratelimit"`
	Webhook    WebhookConfig    `mapstructure:"webhook"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}

type TelegramConfig struct {
	Token   string `mapstructure:"token"`
	AdminID int64  `mapstructure:"admin"`
}

type ModelsConfig struct {
	Timeout           time.Duration `mapstructure:"timeout"`
	IoNetToken        string        `mapstructure:"ioNetToken"`
	PollinationsToken string        `mapstructure:"pollinationsToken"`
	OpenRouterToken   string        `mapstructure:"openRouterToken"`
}

type DatabaseConfig struct {
	Path       string `mapstructure:"path"`
	MaxHistory int    `mapstructure:"maxHistory"`
}

type RatelimitConfig struct {
	Rate int           `mapstructure:"rate"`
	Time time.Duration `mapstructure:"time"`
}

type WebhookConfig struct {
	Port    int    `mapstructure:"port"`
	Enabled bool   `mapstructure:"enabled"`
	Domain  string `mapstructure:"domain"`
}

type PrometheusConfig struct {
	Port    int  `mapstructure:"port"`
	Enabled bool `mapstructure:"enabled"`
}

func LoadConfig() (*Config, error) {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("config")

	config.SetDefault("models.timeout", "2m")

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := config.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Telegram.AdminID == 0 {
		return fmt.Errorf("telegram admin ID is required")
	}

	if c.Telegram.Token == "" {
		return fmt.Errorf("telegram token is required")
	}

	if c.Webhook.Enabled && c.Webhook.Domain == "" {
		return fmt.Errorf("webhook domain is required")
	}

	if c.Models.Timeout == 0 {
		return fmt.Errorf("set correct timeout")
	}

	return nil
}
