package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/toshikimiyagawa/codex-approval-gate/internal/judge"
)

func TestClientSendsChatCompletionRequestAndParsesDecision(t *testing.T) {
	var seenPath string
	var seenAuth string
	var seenModel string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenAuth = r.Header.Get("Authorization")

		var body struct {
			Model    string           `json:"model"`
			Messages []messageForTest `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		seenModel = body.Model
		if len(body.Messages) != 2 {
			t.Fatalf("messages len = %d, want 2", len(body.Messages))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"decision\":\"allow\",\"reason\":\"read-only\"}"}}]}`))
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Model:   "local-model",
		APIKey:  "secret",
		Timeout: time.Second,
	})

	resp, err := client.Decide(context.Background(), judge.ProviderRequest{
		SystemPrompt: "system",
		UserPrompt:   "user",
	})
	if err != nil {
		t.Fatal(err)
	}

	if seenPath != "/v1/chat/completions" {
		t.Fatalf("path = %q, want /v1/chat/completions", seenPath)
	}
	if seenAuth != "Bearer secret" {
		t.Fatalf("authorization = %q, want Bearer secret", seenAuth)
	}
	if seenModel != "local-model" {
		t.Fatalf("model = %q, want local-model", seenModel)
	}
	if resp.Decision != "allow" {
		t.Fatalf("decision = %q, want allow", resp.Decision)
	}
	if resp.Reason != "read-only" {
		t.Fatalf("reason = %q, want read-only", resp.Reason)
	}
}

func TestClientReturnsErrorForNonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, Model: "local-model", Timeout: time.Second})

	_, err := client.Decide(context.Background(), judge.ProviderRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

type messageForTest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
