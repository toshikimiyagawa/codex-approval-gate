package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/toshikimiyagawa/codex-approval-gate/internal/adapters/codex"
	"github.com/toshikimiyagawa/codex-approval-gate/internal/audit"
	"github.com/toshikimiyagawa/codex-approval-gate/internal/config"
	"github.com/toshikimiyagawa/codex-approval-gate/internal/judge"
	"github.com/toshikimiyagawa/codex-approval-gate/internal/policy"
	"github.com/toshikimiyagawa/codex-approval-gate/internal/providers/openai"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: codex-approval-gate codex --config approval-gate.toml")
		return 2
	}

	switch args[0] {
	case "codex":
		return runCodex(args[1:], stdin, stdout, stderr)
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		return 2
	}
}

func runCodex(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("codex", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", "", "path to approval-gate TOML config")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *configPath == "" {
		_, _ = fmt.Fprintln(stderr, "--config is required")
		return 2
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return writeFallback(stdout, codex.OutputModeCodex)
	}

	input, err := io.ReadAll(stdin)
	if err != nil {
		return writeFallback(stdout, cfg.Codex.OutputMode)
	}
	permissionReq, err := codex.DecodePermissionRequest(input)
	if err != nil {
		return writeFallback(stdout, cfg.Codex.OutputMode)
	}

	provider := openai.New(openai.Config{
		BaseURL: cfg.Provider.BaseURL,
		Model:   cfg.Provider.Model,
		APIKey:  apiKey(cfg.Provider.APIKeyEnv),
		Timeout: cfg.Provider.Timeout,
	})
	policyResult := policy.Evaluate(policy.Request{
		ToolName: permissionReq.ToolName,
		Command:  permissionReq.Command,
		CWD:      permissionReq.CWD,
	})
	result := judge.New(provider).Decide(context.Background(), judge.Request{
		ToolName:      permissionReq.ToolName,
		Command:       permissionReq.Command,
		CWD:           permissionReq.CWD,
		Reason:        permissionReq.Reason,
		PolicyVerdict: policyResult.Verdict,
		PolicyReasons: policyResult.Reasons,
		Raw:           permissionReq.Raw,
		Fields:        permissionReq.Fields,
	})

	auditWriter := audit.New(audit.Config{
		Path:            cfg.Audit.Path,
		IncludeRawInput: cfg.Audit.IncludeRawInput,
	})
	if err := auditWriter.Write(audit.Record{
		Decision:     result.Decision,
		Reason:       result.Reason,
		ProviderType: cfg.Provider.Type,
		Model:        cfg.Provider.Model,
		Summary: audit.Summary{
			ToolName: permissionReq.ToolName,
			Command:  permissionReq.Command,
			CWD:      permissionReq.CWD,
		},
		RawInput: permissionReq.Raw,
	}); err != nil {
		result.Decision = judge.DecisionAsk
		result.Reason = "audit write failed"
	}

	output, err := codex.EncodeResponse(result.Decision, cfg.Codex.OutputMode)
	if err != nil {
		return writeFallback(stdout, codex.OutputModeCodex)
	}
	_, _ = stdout.Write(output)
	return 0
}

func writeFallback(stdout io.Writer, outputMode string) int {
	output, err := codex.EncodeResponse(judge.DecisionAsk, outputMode)
	if err != nil {
		output, _ = codex.EncodeResponse(judge.DecisionAsk, codex.OutputModeCodex)
	}
	_, _ = stdout.Write(output)
	return 0
}

func apiKey(envName string) string {
	if envName == "" {
		return ""
	}
	return os.Getenv(envName)
}
