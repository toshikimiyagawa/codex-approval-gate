package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Provider ProviderConfig
	Codex    CodexConfig
	Audit    AuditConfig
	Policy   PolicyConfig
}

type ProviderConfig struct {
	Type      string
	BaseURL   string
	Model     string
	APIKeyEnv string
	Timeout   time.Duration
}

type CodexConfig struct {
	OutputMode string
}

type AuditConfig struct {
	Path            string
	IncludeRawInput bool
}

type PolicyConfig struct {
	ReadOnlyPrefixes []string
	RiskyPrefixes    []string
	RiskySubstrings  []string
}

type fileConfig struct {
	Provider struct {
		Type      string `toml:"type"`
		BaseURL   string `toml:"base_url"`
		Model     string `toml:"model"`
		APIKeyEnv string `toml:"api_key_env"`
		Timeout   string `toml:"timeout"`
	} `toml:"provider"`
	Codex struct {
		OutputMode string `toml:"output_mode"`
	} `toml:"codex"`
	Audit struct {
		Path            string `toml:"path"`
		IncludeRawInput bool   `toml:"include_raw_input"`
	} `toml:"audit"`
	Policy struct {
		ReadOnlyPrefixes []string `toml:"read_only_prefixes"`
		RiskyPrefixes    []string `toml:"risky_prefixes"`
		RiskySubstrings  []string `toml:"risky_substrings"`
	} `toml:"policy"`
}

func Load(path string) (Config, error) {
	var raw fileConfig
	if _, err := toml.DecodeFile(path, &raw); err != nil {
		return Config{}, err
	}

	timeout := 10 * time.Second
	if raw.Provider.Timeout != "" {
		parsed, err := time.ParseDuration(raw.Provider.Timeout)
		if err != nil {
			return Config{}, fmt.Errorf("provider.timeout: %w", err)
		}
		timeout = parsed
	}

	cfg := Config{
		Provider: ProviderConfig{
			Type:      raw.Provider.Type,
			BaseURL:   raw.Provider.BaseURL,
			Model:     raw.Provider.Model,
			APIKeyEnv: raw.Provider.APIKeyEnv,
			Timeout:   timeout,
		},
		Codex: CodexConfig{
			OutputMode: raw.Codex.OutputMode,
		},
		Audit: AuditConfig{
			Path:            raw.Audit.Path,
			IncludeRawInput: raw.Audit.IncludeRawInput,
		},
		Policy: PolicyConfig{
			ReadOnlyPrefixes: raw.Policy.ReadOnlyPrefixes,
			RiskyPrefixes:    raw.Policy.RiskyPrefixes,
			RiskySubstrings:  raw.Policy.RiskySubstrings,
		},
	}
	if cfg.Codex.OutputMode == "" {
		cfg.Codex.OutputMode = "codex"
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Provider.Type == "" {
		return errors.New("provider.type is required")
	}
	if cfg.Provider.BaseURL == "" {
		return errors.New("provider.base_url is required")
	}
	if cfg.Provider.Model == "" {
		return errors.New("provider.model is required")
	}
	if cfg.Codex.OutputMode != "codex" && cfg.Codex.OutputMode != "simple" {
		return errors.New("codex.output_mode must be codex or simple")
	}
	return nil
}
