package config

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/config"
)

const (
	LocalPath = "./config/config.yaml"
)

type Config struct {
	Api      ApiCfg      `mapstructure:"api" yaml:"api"`
	RabbitMQ RabbitMQCfg `mapstructure:"rabbitmq" yaml:"rabbitmq"`
	Telegram TelegramCfg `mapstructure:"telegram" yaml:"telegram"`
}

type ApiCfg struct {
	GinMode      string `mapstructure:"gin_mode" yaml:"gin_mode"`
	Addr         string `mapstructure:"addr" yaml:"addr"`
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

type RabbitMQCfg struct {
	URL        string `mapstructure:"url" yaml:"url"`
	Exchange   string `mapstructure:"exchange" yaml:"exchange"`
	Queue      string `mapstructure:"queue" yaml:"queue"`
	RoutingKey string `mapstructure:"routing_key" yaml:"routing_key"`
}

type TelegramCfg struct {
	BotToken      string `mapstructure:"bot_token" yaml:"bot_token"`
	DefaultChatId string `mapstructure:"default_chat_id" yaml:"default_chat_id"`
}

func (cfg *Config) Read(paths ...string) error {
	c := config.New()

	if err := c.LoadConfigFiles(paths...); err != nil {
		return fmt.Errorf("{Read 1}: %w", err)
	}

	if err := c.Unmarshal(cfg); err != nil {
		return fmt.Errorf("{Read 2}: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return fmt.Errorf("{Read 3}: %w", err)
	}

	cfg.print()
	return nil
}

func (cfg *Config) print() {
	log.Println("=================== CONFIG ===================")
	log.Println("Gin Mode...............", cfg.Api.GinMode)
	log.Println("API Address............", cfg.Api.Addr)
	log.Println("Read Timeout...........", cfg.Api.ReadTimeout)
	log.Println("Write Timeout..........", cfg.Api.WriteTimeout)
	log.Println("Idle Timeout...........", cfg.Api.IdleTimeout)
	log.Println("RabbitMQ URL...........", cfg.RabbitMQ.URL)
	log.Println("RabbitMQ Exchange......", cfg.RabbitMQ.Exchange)
	log.Println("RabbitMQ Queue.........", cfg.RabbitMQ.Queue)
	log.Println("RabbitMQ RoutingKey....", cfg.RabbitMQ.RoutingKey)
	log.Println("Telegram ChatID........", cfg.Telegram.DefaultChatId)
	log.Printf("==============================================\n\n")
}

func (cfg *Config) validate() error {
	if cfg.Api.GinMode == "" || (cfg.Api.GinMode != gin.DebugMode && cfg.Api.GinMode != gin.ReleaseMode) {
		return fmt.Errorf("invalid gin mode: %s", cfg.Api.GinMode)
	}

	if cfg.Api.Addr == "" {
		return fmt.Errorf("invalid addr: %s", cfg.Api.Addr)
	}

	if cfg.Api.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout: %d", cfg.Api.ReadTimeout)
	}

	if cfg.Api.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout: %d", cfg.Api.WriteTimeout)
	}

	if cfg.Api.IdleTimeout <= 0 {
		return fmt.Errorf("invalid idle timeout: %d", cfg.Api.IdleTimeout)
	}

	if cfg.RabbitMQ.URL == "" {
		return fmt.Errorf("rabbitmq url is required")
	}
	if cfg.RabbitMQ.Exchange == "" {
		return fmt.Errorf("rabbitmq exchange is required")
	}
	if cfg.RabbitMQ.Queue == "" {
		return fmt.Errorf("rabbitmq queue is required")
	}
	if cfg.RabbitMQ.RoutingKey == "" {
		return fmt.Errorf("rabbitmq routing_key is required")
	}

	return nil
}
