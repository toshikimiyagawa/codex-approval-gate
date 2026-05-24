package judge

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestJudgeReturnsProviderDecision(t *testing.T) {
	j := New(staticProvider{response: ProviderResponse{Decision: "allow", Reason: "safe read"}})

	result := j.Decide(context.Background(), Request{ToolName: "shell", Command: "git status"})

	if result.Decision != "allow" {
		t.Fatalf("decision = %q, want allow", result.Decision)
	}
	if result.Reason != "safe read" {
		t.Fatalf("reason = %q, want safe read", result.Reason)
	}
}

func TestJudgeFallsBackToAskOnProviderError(t *testing.T) {
	j := New(staticProvider{err: errors.New("offline")})

	result := j.Decide(context.Background(), Request{Command: "rm -rf /tmp/x"})

	if result.Decision != "ask" {
		t.Fatalf("decision = %q, want ask", result.Decision)
	}
}

func TestJudgeFallsBackToAskOnUnknownDecision(t *testing.T) {
	j := New(staticProvider{response: ProviderResponse{Decision: "maybe"}})

	result := j.Decide(context.Background(), Request{Command: "go test ./..."})

	if result.Decision != "ask" {
		t.Fatalf("decision = %q, want ask", result.Decision)
	}
}

func TestJudgeBuildsPromptFromRequest(t *testing.T) {
	provider := &recordingProvider{response: ProviderResponse{Decision: "deny"}}
	j := New(provider)

	result := j.Decide(context.Background(), Request{
		ToolName: "shell",
		Command:  "curl http://example.com",
		CWD:      "/tmp/project",
		Reason:   "fetch data",
	})

	if result.Decision != "deny" {
		t.Fatalf("decision = %q, want deny", result.Decision)
	}
	if provider.request.Command != "curl http://example.com" {
		t.Fatalf("provider command = %q, want curl http://example.com", provider.request.Command)
	}
	if provider.request.SystemPrompt == "" {
		t.Fatal("system prompt was empty")
	}
	if provider.request.UserPrompt == "" {
		t.Fatal("user prompt was empty")
	}
}

func TestJudgeIncludesPolicyResultInPrompt(t *testing.T) {
	provider := &recordingProvider{response: ProviderResponse{Decision: "ask"}}
	j := New(provider)

	j.Decide(context.Background(), Request{
		ToolName:      "shell",
		Command:       "rm -rf /tmp/build",
		PolicyVerdict: "risky",
		PolicyReasons: []string{"starts with risky command rm"},
	})

	if !strings.Contains(provider.request.UserPrompt, `"policy_verdict":"risky"`) {
		t.Fatalf("user prompt = %s, want policy verdict", provider.request.UserPrompt)
	}
	if !strings.Contains(provider.request.UserPrompt, "starts with risky command rm") {
		t.Fatalf("user prompt = %s, want policy reason", provider.request.UserPrompt)
	}
}

type staticProvider struct {
	response ProviderResponse
	err      error
}

func (p staticProvider) Decide(context.Context, ProviderRequest) (ProviderResponse, error) {
	return p.response, p.err
}

type recordingProvider struct {
	response ProviderResponse
	request  ProviderRequest
}

func (p *recordingProvider) Decide(_ context.Context, req ProviderRequest) (ProviderResponse, error) {
	p.request = req
	return p.response, nil
}
