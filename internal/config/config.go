package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all APSON configuration options.
type Config struct {
	Email struct {
		Sender     string   `yaml:"sender"`
		Password   string   `yaml:"password"`
		Recipients []string `yaml:"recipients"`
		SMTPServer string   `yaml:"smtp_server"`
		SMTPPort   int      `yaml:"smtp_port"`
	} `yaml:"email"`
	PollingIntervalMinutes int      `yaml:"polling_interval_minutes"`
	NotifyWindowDays       int      `yaml:"notify_window_days"`
	Buildings              []string `yaml:"buildings"`
}

// LoadConfig loads configuration from the given YAML file path.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}

	// Set defaults if not present
	if cfg.PollingIntervalMinutes == 0 {
		cfg.PollingIntervalMinutes = 60
	}
	if cfg.NotifyWindowDays == 0 {
		cfg.NotifyWindowDays = 14
	}
	if len(cfg.Buildings) == 0 {
		cfg.Buildings = []string{"CPH"}
	}

	return &cfg, nil
}
