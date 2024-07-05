package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	DatabaseURL string `yaml:"database_url" env-required:"true"`
	HttpServerConfig
}

type HttpServerConfig struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func MustLoadEnv() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
}
