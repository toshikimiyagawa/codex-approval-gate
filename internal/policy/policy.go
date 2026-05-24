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

type Result struct {
	Verdict string   `json:"verdict"`
	Reasons []string `json:"reasons,omitempty"`
}

func Evaluate(req Request) Result {
	command := strings.TrimSpace(req.Command)
	if command == "" {
		return Result{Verdict: VerdictUnknown}
	}

	if reasons := riskyReasons(command); len(reasons) > 0 {
		return Result{Verdict: VerdictRisky, Reasons: reasons}
	}
	if isReadOnly(command) {
		return Result{Verdict: VerdictSafe, Reasons: []string{"matches a read-only command pattern"}}
	}
	return Result{Verdict: VerdictUnknown}
}

func riskyReasons(command string) []string {
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
	return reasons
}

func isReadOnly(command string) bool {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false
	}

	switch fields[0] {
	case "pwd", "ls", "cat":
		return true
	case "git":
		return len(fields) >= 2 && (fields[1] == "status" || fields[1] == "diff")
	default:
		return false
	}
}
