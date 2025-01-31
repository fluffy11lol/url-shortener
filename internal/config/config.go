package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server" env-required:"true"`
}
type HttpServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
		//return nil, errors.New("CONFIG_PATH is not set")
	}

	config := &Config{}

	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		return nil, err
	}

	return config, nil
}
