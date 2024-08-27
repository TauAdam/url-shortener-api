package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env              string `yaml:"env" env-required:"true"`
	DatabaseURL      string `yaml:"database_url" env-required:"true"`
	HttpServerConfig `yaml:"http_server_config"`
	JWTSecret        string      `yaml:"jwt_secret" env-required:"true" env:"JWT_SECRET"`
	Clients          ClientsList `yaml:"clients"`
}

type HttpServerConfig struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"SERVER_PASSWORD"`
}

type Client struct {
	Address       string        `yaml:"address"`
	Timeout       time.Duration `yaml:"timeout"`
	RetriesNumber int           `yaml:"retries_number"`
	IsSecure      bool          `yaml:"is_secure"`
}

type ClientsList struct {
	SSO Client `yaml:"sso"`
}

func MustLoadEnv() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}
	return &cfg
}
