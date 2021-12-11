package api

import (
	"forum/internal/storage"
)

type Config struct {
	Port    string
	Storage *storage.Config
}

func NewConfig(port string, db string) *Config {
	return &Config{
		Port:    port,
		Storage: storage.NewConfig(db),
	}
}
