package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCodexReturnsAllowResponse(t *testing.T) {
	server := providerServer(t, `{"decision":"allow","reason":"read-only"}`, http.StatusOK)
	defer server.Close()
	cfg := writeCLIConfig(t, server.URL, "codex", "")

	var stdout bytes.Buffer
	code := run([]string{"codex", "--config", cfg}, strings.NewReader(hookInput), &stdout, &bytes.Buffer{})

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	var decoded struct {
		HookSpecificOutput struct {
			PermissionDecision string `json:"permissionDecision"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.HookSpecificOutput.PermissionDecision != "allow" {
		t.Fatalf("decision = %q, want allow", decoded.HookSpecificOutput.PermissionDecision)
	}
}

func TestRunCodexFallsBackToAskOnProviderError(t *testing.T) {
	server := providerServer(t, `nope`, http.StatusInternalServerError)
	defer server.Close()
	cfg := writeCLIConfig(t, server.URL, "simple", "")

	var stdout bytes.Buffer
	code := run([]string{"codex", "--config", cfg}, strings.NewReader(hookInput), &stdout, &bytes.Buffer{})

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if strings.TrimSpace(stdout.String()) != `{"decision":"ask"}` {
		t.Fatalf("stdout = %q, want simple ask", stdout.String())
	}
}

func TestRunCodexDiscoversLocalConfigWhenConfigFlagOmitted(t *testing.T) {
	server := providerServer(t, `{"decision":"allow","reason":"local config"}`, http.StatusOK)
	defer server.Close()
	tmpDir := t.TempDir()
	writeCLIConfigAt(t, filepath.Join(tmpDir, "approval-gate.toml"), server.URL, "simple", "")
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	code := run([]string{"codex"}, strings.NewReader(hookInput), &stdout, &bytes.Buffer{})

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if strings.TrimSpace(stdout.String()) != `{"decision":"allow"}` {
		t.Fatalf("stdout = %q, want simple allow", stdout.String())
	}
}

func TestResolveConfigPathDiscoversUserConfigDir(t *testing.T) {
	configHome := t.TempDir()
	t.Setenv("CODEX_APPROVAL_GATE_CONFIG_HOME", configHome)
	configPath := filepath.Join(configHome, "codex-approval-gate", "config.toml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := resolveConfigPath("")
	if err != nil {
		t.Fatal(err)
	}
	if got != configPath {
		t.Fatalf("config path = %q, want %q", got, configPath)
	}
}

func TestRunCodexWritesAuditRecord(t *testing.T) {
	server := providerServer(t, `{"decision":"deny","reason":"dangerous"}`, http.StatusOK)
	defer server.Close()
	auditPath := filepath.Join(t.TempDir(), "audit.jsonl")
	cfg := writeCLIConfig(t, server.URL, "simple", auditPath)

	var stdout bytes.Buffer
	code := run([]string{"codex", "--config", cfg}, strings.NewReader(hookInput), &stdout, &bytes.Buffer{})

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	contents, err := os.ReadFile(auditPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), `"decision":"deny"`) {
		t.Fatalf("audit = %s, want deny decision", contents)
	}
}

func TestRunCodexIncludesPolicyPrecheckInProviderPrompt(t *testing.T) {
	var userPrompt string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Messages []messageForTest `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		for _, message := range body.Messages {
			if message.Role == "user" {
				userPrompt = message.Content
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"decision\":\"ask\"}"}}]}`))
	}))
	defer server.Close()
	cfg := writeCLIConfig(t, server.URL, "simple", "")

	var stdout bytes.Buffer
	code := run([]string{"codex", "--config", cfg}, strings.NewReader(`{
  "hook_event_name": "PermissionRequest",
  "tool_name": "shell",
  "cwd": "/tmp/project",
  "command": "rm -rf /tmp/build"
}`), &stdout, &bytes.Buffer{})

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(userPrompt, `"policy_verdict":"risky"`) {
		t.Fatalf("user prompt = %s, want risky policy verdict", userPrompt)
	}
	if !strings.Contains(userPrompt, "starts with risky command rm") {
		t.Fatalf("user prompt = %s, want risky policy reason", userPrompt)
	}
}

func providerServer(t *testing.T, content string, status int) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q, want /v1/chat/completions", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status >= 200 && status < 300 {
			_, _ = w.Write([]byte(`{"choices":[{"message":{"content":` + strconvQuote(content) + `}}]}`))
			return
		}
		_, _ = w.Write([]byte(content))
	}))
}

func writeCLIConfig(t *testing.T, baseURL string, outputMode string, auditPath string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "approval-gate.toml")
	writeCLIConfigAt(t, path, baseURL, outputMode, auditPath)
	return path
}

func writeCLIConfigAt(t *testing.T, path string, baseURL string, outputMode string, auditPath string) {
	t.Helper()

	contents := `[provider]
type = "openai"
base_url = "` + baseURL + `"
model = "local-model"
timeout = "3s"

[codex]
output_mode = "` + outputMode + `"

[audit]
path = "` + auditPath + `"
`
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func strconvQuote(value string) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

type messageForTest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const hookInput = `{
  "hook_event_name": "PermissionRequest",
  "tool_name": "shell",
  "cwd": "/tmp/project",
  "command": "git status --short",
  "reason": "inspect"
}`
