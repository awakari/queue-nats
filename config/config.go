package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Api struct {
		Port uint16 `envconfig:"API_PORT" default:"8080" required:"true"`
	}
	Log struct {
		Level int `envconfig:"LOG_LEVEL" default:"-4" required:"true"`
	}
	Nats NatsConfig
}

type NatsConfig struct {
	Uri               string `envconfig:"NATS_URI" default:"nats://nats:4222" required:"true"`
	User              string `envconfig:"NATS_USER" default:"awakari" required:"true"`
	Password          string `envconfig:"NATS_PASSWORD" default:"awakari" required:"true"`
	PollTimeoutMillis uint32 `envconfig:"NATS_POLL_TIMEOUT_MILLIS" default:"1" required:"true"`
}

func NewConfigFromEnv() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
