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
	return path
}

func strconvQuote(value string) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

const hookInput = `{
  "hook_event_name": "PermissionRequest",
  "tool_name": "shell",
  "cwd": "/tmp/project",
  "command": "git status --short",
  "reason": "inspect"
}`
