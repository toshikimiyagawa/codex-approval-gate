package judge

import (
	"context"
	"encoding/json"
	"strings"
)

const (
	DecisionAllow = "allow"
	DecisionDeny  = "deny"
	DecisionAsk   = "ask"
)

type Request struct {
	ToolName      string         `json:"tool_name,omitempty"`
	Command       string         `json:"command,omitempty"`
	CWD           string         `json:"cwd,omitempty"`
	Reason        string         `json:"reason,omitempty"`
	PolicyVerdict string         `json:"policy_verdict,omitempty"`
	PolicyReasons []string       `json:"policy_reasons,omitempty"`
	Raw           []byte         `json:"-"`
	Fields        map[string]any `json:"fields,omitempty"`
}

type Result struct {
	Decision string
	Reason   string
}

type Provider interface {
	Decide(context.Context, ProviderRequest) (ProviderResponse, error)
}

type ProviderRequest struct {
	ToolName     string
	Command      string
	CWD          string
	Reason       string
	SystemPrompt string
	UserPrompt   string
}

type ProviderResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
}

type Judge struct {
	provider Provider
}

func New(provider Provider) Judge {
	return Judge{provider: provider}
}

func (j Judge) Decide(ctx context.Context, req Request) Result {
	providerReq := ProviderRequest{
		ToolName:     req.ToolName,
		Command:      req.Command,
		CWD:          req.CWD,
		Reason:       req.Reason,
		SystemPrompt: systemPrompt(),
		UserPrompt:   userPrompt(req),
	}

	response, err := j.provider.Decide(ctx, providerReq)
	if err != nil {
		return Result{Decision: DecisionAsk, Reason: "provider error"}
	}
	decision := strings.ToLower(strings.TrimSpace(response.Decision))
	if !validDecision(decision) {
		return Result{Decision: DecisionAsk, Reason: "invalid provider decision"}
	}
	return Result{Decision: decision, Reason: response.Reason}
}

func ParseProviderResponse(content string) (ProviderResponse, error) {
	var response ProviderResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return ProviderResponse{}, err
	}
	return response, nil
}

func systemPrompt() string {
	return "You are an approval gate for Codex PermissionRequest hooks. Return JSON only with decision allow, deny, or ask and an optional reason. Return ask if unsure."
}

func userPrompt(req Request) string {
	payload := map[string]any{
		"tool_name":      req.ToolName,
		"command":        req.Command,
		"cwd":            req.CWD,
		"reason":         req.Reason,
		"policy_verdict": req.PolicyVerdict,
		"policy_reasons": req.PolicyReasons,
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

func validDecision(decision string) bool {
	return decision == DecisionAllow || decision == DecisionDeny || decision == DecisionAsk
}
