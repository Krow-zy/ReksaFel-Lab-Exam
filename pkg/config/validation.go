package config

import (
	"errors"
	"net"
)

// Config represents server parameters to validate.
type Config struct {
	Port         int      `json:"port"`
	SecretKey    string   `json:"secret_key"`
	IPWhitelist  []string `json:"ip_whitelist"`
}

// ValidateConfig verifies if the server config properties are safe.
// Uses standard library net package for IP validation to remain lightweight.
func ValidateConfig(c *Config) error {
	if c.Port < 1024 || c.Port > 65535 {
		return errors.New("port must be in the user range 1024-65535")
	}

	if len(c.SecretKey) < 16 {
		return errors.New("secret key must be at least 16 characters long")
	}

	for _, ip := range c.IPWhitelist {
		if net.ParseIP(ip) == nil {
			return errors.New("invalid IP address in whitelist: " + ip)
		}
	}

	return nil
}
