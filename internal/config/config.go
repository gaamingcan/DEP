package config

import (
	"fmt"
	"os"
	"sort"

	"github.com/pelletier/go-toml/v2"
)

const FileName = "dep.toml"

// Repository describes a Git repository in dep.toml.
type Repository struct {
	URL string `toml:"url"`
}

// Config is the DEP project configuration.
type Config struct {
	Repositories []Repository `toml:"repo"`
}

// Read parses dep.toml from path.
func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &cfg, nil
}

// Write writes dep.toml to path with repositories sorted by URL.
func Write(path string, cfg *Config) error {
	sortRepositories(cfg)
	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// ContainsURL reports whether url is already configured.
func (cfg *Config) ContainsURL(url string) bool {
	for _, r := range cfg.Repositories {
		if r.URL == url {
			return true
		}
	}
	return false
}

// Add appends a repository URL if not already present.
func (cfg *Config) Add(url string) error {
	if cfg.ContainsURL(url) {
		return fmt.Errorf("repository URL already exists: %s", url)
	}
	cfg.Repositories = append(cfg.Repositories, Repository{URL: url})
	sortRepositories(cfg)
	return nil
}

// Remove deletes a repository by URL.
func (cfg *Config) Remove(url string) error {
	for i, r := range cfg.Repositories {
		if r.URL == url {
			cfg.Repositories = append(cfg.Repositories[:i], cfg.Repositories[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("repository not found: %s", url)
}

// NewEmpty returns an empty project configuration.
func NewEmpty() *Config {
	return &Config{Repositories: []Repository{}}
}

func sortRepositories(cfg *Config) {
	sort.Slice(cfg.Repositories, func(i, j int) bool {
		return cfg.Repositories[i].URL < cfg.Repositories[j].URL
	})
}

// URLs returns all configured repository URLs.
func (cfg *Config) URLs() []string {
	urls := make([]string, len(cfg.Repositories))
	for i, r := range cfg.Repositories {
		urls[i] = r.URL
	}
	return urls
}
