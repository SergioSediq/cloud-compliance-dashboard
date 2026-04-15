package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Config holds process-level settings loaded from the environment.
type Config struct {
	ChecksPath string
	Port       string
}

// Load reads CHECKS_PATH and PORT with the same defaults as the original single-binary layout.
func Load() Config {
	p := os.Getenv("CHECKS_PATH")
	if p == "" {
		p = filepath.Join("data", "checks.json")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	return Config{ChecksPath: p, Port: port}
}

// ListenAddr returns a host:port form suitable for net/http.
func (c Config) ListenAddr() string {
	if strings.HasPrefix(c.Port, ":") {
		return c.Port
	}
	return ":" + c.Port
}
