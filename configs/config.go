package configs

import (
	"log"
	"os"
)

const (
	defaultTTHost     = "localhost"
	defaultTTPort     = "3301"
	defaultServerPort = "80"

	envTTHost     = "TARANTOOL_HOST"
	envTTPort     = "TARANTOOL_PORT"
	envServerPort = "SERVER_PORT"
)

type TTStoreConfig struct {
	Host string
	Port string
}

type TTOption func(*TTStoreConfig)

func WithTTAddress(host, port string) TTOption {
	return func(cfg *TTStoreConfig) {
		cfg.Host = host
		cfg.Port = port
	}
}

func NewTTStoreConfig(opts ...TTOption) TTStoreConfig {
	cfg := TTStoreConfig{
		Host: os.Getenv(envTTHost),
		Port: os.Getenv(envTTPort),
	}

	if cfg.Host == "" {
		log.Println("loading default host for tarantool")
		cfg.Host = defaultTTHost
	}
	if cfg.Port == "" {
		log.Println("loading default port for tarantool")
		cfg.Port = defaultTTPort
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

type ServerConfig struct {
	Port string
}

type ServerOption func(*ServerConfig)

func WithPort(port string) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.Port = port
	}
}

func NewServerConfig(opts ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		Port: os.Getenv(envServerPort),
	}

	if cfg.Port == "" {
		log.Println("loading default port")
		cfg.Port = defaultServerPort
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

type Config struct {
	TT     TTStoreConfig
	Server ServerConfig
}

func Load() *Config {
	return &Config{
		TT:     NewTTStoreConfig(),
		Server: *NewServerConfig(),
	}
}
