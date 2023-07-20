package config

import (
	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	DatabaseConfig
	ClientConfig
}

type DatabaseConfig struct {
	DSN string `envconfig:"POSTGRES_DSN" required:"true"`
}

type ClientConfig struct {
	NotificationServiceURL string `envconfig:"NOTIFICATION_URL" required:"true"`
}

// Load reads the configuration from environment variables
func (c *Config) Load() {
	log.Infof("loading configuration from environment variables")
	err := envconfig.Process("COMPUTER_MANAGER_API", c)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
}
