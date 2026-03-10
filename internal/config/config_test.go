package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTestConfig(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	orig := os.Getenv("HOME")
	// Override HOME so configDir() uses our temp dir
	t.Setenv("HOME", dir)
	return dir, func() {
		os.Setenv("HOME", orig)
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{
		BaseURL:     "https://analytics.example.com",
		DefaultSite: "example.com",
		Sites: map[string]Site{
			"example.com": {Token: "tok123"},
		},
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.BaseURL != cfg.BaseURL {
		t.Errorf("BaseURL = %q, want %q", loaded.BaseURL, cfg.BaseURL)
	}
	if loaded.DefaultSite != cfg.DefaultSite {
		t.Errorf("DefaultSite = %q, want %q", loaded.DefaultSite, cfg.DefaultSite)
	}
	site, ok := loaded.Sites["example.com"]
	if !ok {
		t.Fatal("site example.com not found")
	}
	if site.Token != "tok123" {
		t.Errorf("Token = %q, want %q", site.Token, "tok123")
	}
}

func TestLoadNoFile(t *testing.T) {
	setupTestConfig(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.BaseURL != "https://plausible.io" {
		t.Errorf("default BaseURL = %q, want %q", cfg.BaseURL, "https://plausible.io")
	}
	if len(cfg.Sites) != 0 {
		t.Errorf("Sites should be empty, got %d", len(cfg.Sites))
	}
}

func TestAddSite(t *testing.T) {
	setupTestConfig(t)

	// Save initial config
	if err := Save(&Config{BaseURL: "https://plausible.io", Sites: make(map[string]Site)}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := AddSite("blog.example.com", "tok-blog"); err != nil {
		t.Fatalf("AddSite: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	site, ok := cfg.Sites["blog.example.com"]
	if !ok {
		t.Fatal("site not found after AddSite")
	}
	if site.Token != "tok-blog" {
		t.Errorf("Token = %q, want %q", site.Token, "tok-blog")
	}
	// First site added should become default
	if cfg.DefaultSite != "blog.example.com" {
		t.Errorf("DefaultSite = %q, want %q", cfg.DefaultSite, "blog.example.com")
	}
}

func TestRemoveSite(t *testing.T) {
	setupTestConfig(t)

	if err := Save(&Config{
		BaseURL:     "https://plausible.io",
		DefaultSite: "example.com",
		Sites:       map[string]Site{"example.com": {Token: "tok"}},
	}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := RemoveSite("example.com"); err != nil {
		t.Fatalf("RemoveSite: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if _, ok := cfg.Sites["example.com"]; ok {
		t.Error("site should be removed")
	}
	if cfg.DefaultSite != "" {
		t.Errorf("DefaultSite should be cleared, got %q", cfg.DefaultSite)
	}
}

func TestRemoveSiteNotFound(t *testing.T) {
	setupTestConfig(t)

	if err := Save(&Config{BaseURL: "https://plausible.io", Sites: make(map[string]Site)}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	err := RemoveSite("nonexistent.com")
	if err == nil {
		t.Fatal("expected error removing nonexistent site")
	}
}

func TestGetSite(t *testing.T) {
	dir, _ := setupTestConfig(t)

	// Write config file directly
	cfgDir := filepath.Join(dir, ".config", "plaus")
	os.MkdirAll(cfgDir, 0755)
	data, _ := json.Marshal(Config{
		BaseURL: "https://plausible.io",
		Sites:   map[string]Site{"test.com": {Token: "abc"}},
	})
	os.WriteFile(filepath.Join(cfgDir, "config.json"), data, 0644)

	site, ok := GetSite("test.com")
	if !ok {
		t.Fatal("GetSite should find test.com")
	}
	if site.Token != "abc" {
		t.Errorf("Token = %q, want %q", site.Token, "abc")
	}

	_, ok = GetSite("missing.com")
	if ok {
		t.Error("GetSite should not find missing.com")
	}
}
