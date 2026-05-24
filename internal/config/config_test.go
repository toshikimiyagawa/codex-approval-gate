package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadValidConfigAppliesDefaults(t *testing.T) {
	path := writeConfig(t, `
[provider]
type = "openai"
base_url = "http://localhost:1234"
model = "local-model"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Provider.Type != "openai" {
		t.Fatalf("provider type = %q, want openai", cfg.Provider.Type)
	}
	if cfg.Provider.Timeout != 10*time.Second {
		t.Fatalf("timeout = %s, want 10s", cfg.Provider.Timeout)
	}
	if cfg.Codex.OutputMode != "codex" {
		t.Fatalf("output mode = %q, want codex", cfg.Codex.OutputMode)
	}
}

func TestLoadParsesConfiguredFields(t *testing.T) {
	path := writeConfig(t, `
[provider]
type = "openai"
base_url = "http://localhost:1234"
model = "local-model"
api_key_env = "LOCAL_API_KEY"
timeout = "3s"

[codex]
output_mode = "simple"

[audit]
path = "./audit.jsonl"
include_raw_input = true

[policy]
read_only_prefixes = ["go test"]
risky_prefixes = ["docker system prune"]
risky_substrings = ["| pbcopy"]
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Provider.APIKeyEnv != "LOCAL_API_KEY" {
		t.Fatalf("api_key_env = %q, want LOCAL_API_KEY", cfg.Provider.APIKeyEnv)
	}
	if cfg.Provider.Timeout != 3*time.Second {
		t.Fatalf("timeout = %s, want 3s", cfg.Provider.Timeout)
	}
	if cfg.Codex.OutputMode != "simple" {
		t.Fatalf("output mode = %q, want simple", cfg.Codex.OutputMode)
	}
	if cfg.Audit.Path != "./audit.jsonl" {
		t.Fatalf("audit path = %q, want ./audit.jsonl", cfg.Audit.Path)
	}
	if !cfg.Audit.IncludeRawInput {
		t.Fatal("include raw input = false, want true")
	}
	if got := cfg.Policy.ReadOnlyPrefixes; len(got) != 1 || got[0] != "go test" {
		t.Fatalf("read_only_prefixes = %v, want [go test]", got)
	}
	if got := cfg.Policy.RiskyPrefixes; len(got) != 1 || got[0] != "docker system prune" {
		t.Fatalf("risky_prefixes = %v, want [docker system prune]", got)
	}
	if got := cfg.Policy.RiskySubstrings; len(got) != 1 || got[0] != "| pbcopy" {
		t.Fatalf("risky_substrings = %v, want [| pbcopy]", got)
	}
}

func TestLoadRequiresProviderFields(t *testing.T) {
	path := writeConfig(t, `
[provider]
type = "openai"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "provider.base_url") {
		t.Fatalf("error = %q, want provider.base_url", err.Error())
	}
}

func TestLoadRejectsInvalidTimeout(t *testing.T) {
	path := writeConfig(t, `
[provider]
type = "openai"
base_url = "http://localhost:1234"
model = "local-model"
timeout = "soon"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func writeConfig(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "approval-gate.toml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
