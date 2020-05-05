package config

import (
	"context"

	"github.com/heetch/confita"

	"github.com/IgorRybak2055/bamboo/internal/bamboo"
	"github.com/IgorRybak2055/bamboo/internal/storage"
)

// Config stores configs for ragger
type Config struct {
	HTTP *bamboo.HTTPConfig
	DB   *storage.Config
}

// NewConfig returns Config with values from environment variables
func NewConfig(ctx context.Context) (*Config, error) {
	var cfg = &Config{
		HTTP: &bamboo.HTTPConfig{},
		DB:   &storage.Config{},
	}

	if err := confita.NewLoader().Load(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
