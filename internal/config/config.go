package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	BaseURL     string          `json:"base_url"`
	DefaultSite string          `json:"default_site"`
	Sites       map[string]Site `json:"sites"`
	GraphType   string          `json:"graph_type,omitempty"`
	Period      string          `json:"period,omitempty"`
}

type Site struct {
	Token string `json:"token"`
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "plaus"), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				BaseURL: "https://plausible.io",
				Sites:   make(map[string]Site),
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Sites == nil {
		cfg.Sites = make(map[string]Site)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddSite(domain string, token string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.Sites[domain] = Site{Token: token}
	if cfg.DefaultSite == "" {
		cfg.DefaultSite = domain
	}
	return Save(cfg)
}

func RemoveSite(domain string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if _, ok := cfg.Sites[domain]; !ok {
		return fmt.Errorf("site %q not configured", domain)
	}
	delete(cfg.Sites, domain)
	if cfg.DefaultSite == domain {
		cfg.DefaultSite = ""
	}
	return Save(cfg)
}

func GetSite(domain string) (Site, bool) {
	cfg, err := Load()
	if err != nil {
		return Site{}, false
	}
	site, ok := cfg.Sites[domain]
	return site, ok
}
