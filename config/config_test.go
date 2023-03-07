package config

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := NewConfigFromEnv()
	assert.Nil(t, err)
	assert.Equal(t, uint16(8080), cfg.Api.Port)
	assert.Equal(t, slog.LevelDebug, slog.Level(cfg.Log.Level))
	assert.Equal(t, "nats://nats:4222", cfg.Nats.Uri)
}
