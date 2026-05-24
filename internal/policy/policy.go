package policy

import (
	"strings"
)

const (
	VerdictSafe    = "safe"
	VerdictRisky   = "risky"
	VerdictUnknown = "unknown"
)

type Request struct {
	ToolName string
	Command  string
	CWD      string
}

type Config struct {
	ReadOnlyPrefixes []string
	RiskyPrefixes    []string
	RiskySubstrings  []string
}

type Result struct {
	Verdict string   `json:"verdict"`
	Reasons []string `json:"reasons,omitempty"`
}

func Evaluate(req Request) Result {
	return EvaluateWithConfig(req, Config{})
}

func EvaluateWithConfig(req Request, cfg Config) Result {
	command := strings.TrimSpace(req.Command)
	if command == "" {
		return Result{Verdict: VerdictUnknown}
	}

	if reasons := riskyReasons(command, cfg); len(reasons) > 0 {
		return Result{Verdict: VerdictRisky, Reasons: reasons}
	}
	if isReadOnly(command, cfg) {
		return Result{Verdict: VerdictSafe, Reasons: []string{"matches a read-only command pattern"}}
	}
	return Result{Verdict: VerdictUnknown}
}

func riskyReasons(command string, cfg Config) []string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return nil
	}

	var reasons []string
	switch fields[0] {
	case "rm", "sudo", "chmod", "chown":
		reasons = append(reasons, "starts with risky command "+fields[0])
	}
	if strings.Contains(command, ">") {
		reasons = append(reasons, "contains shell redirection")
	}
	for _, prefix := range cfg.RiskyPrefixes {
		if commandMatchesPrefix(command, prefix) {
			reasons = append(reasons, "matches configured risky prefix "+prefix)
		}
	}
	for _, substring := range cfg.RiskySubstrings {
		if substring != "" && strings.Contains(command, substring) {
			reasons = append(reasons, "contains configured risky substring "+substring)
		}
	}
	return reasons
}

func isReadOnly(command string, cfg Config) bool {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false
	}

	switch fields[0] {
	case "pwd", "ls", "cat":
		return true
	case "git":
		if len(fields) >= 2 && (fields[1] == "status" || fields[1] == "diff") {
			return true
		}
	}
	for _, prefix := range cfg.ReadOnlyPrefixes {
		if commandMatchesPrefix(command, prefix) {
			return true
		}
	}
	return false
}

func commandMatchesPrefix(command string, prefix string) bool {
	prefix = strings.TrimSpace(prefix)
	return prefix != "" && (command == prefix || strings.HasPrefix(command, prefix+" "))
}
