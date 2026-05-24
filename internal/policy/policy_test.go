package policy

import "testing"

func TestEvaluateMarksReadOnlyCommandsSafe(t *testing.T) {
	for _, command := range []string{
		"git status --short",
		"git diff -- README.md",
		"ls -la",
		"pwd",
		"cat README.md",
	} {
		t.Run(command, func(t *testing.T) {
			result := Evaluate(Request{ToolName: "shell", Command: command})
			if result.Verdict != VerdictSafe {
				t.Fatalf("verdict = %q, want %q; reasons: %v", result.Verdict, VerdictSafe, result.Reasons)
			}
		})
	}
}

func TestEvaluateMarksRiskyCommandsRisky(t *testing.T) {
	for _, command := range []string{
		"rm -rf /tmp/build",
		"sudo make install",
		"chmod 600 secret.txt",
		"chown root:wheel file",
		"echo hello > README.md",
	} {
		t.Run(command, func(t *testing.T) {
			result := Evaluate(Request{ToolName: "shell", Command: command})
			if result.Verdict != VerdictRisky {
				t.Fatalf("verdict = %q, want %q; reasons: %v", result.Verdict, VerdictRisky, result.Reasons)
			}
			if len(result.Reasons) == 0 {
				t.Fatal("expected at least one reason")
			}
		})
	}
}

func TestEvaluateMarksUnknownCommandsUnknown(t *testing.T) {
	result := Evaluate(Request{ToolName: "shell", Command: "go test ./..."})
	if result.Verdict != VerdictUnknown {
		t.Fatalf("verdict = %q, want %q", result.Verdict, VerdictUnknown)
	}
}

func TestEvaluateWithConfigAddsReadOnlyPrefixes(t *testing.T) {
	result := EvaluateWithConfig(Request{ToolName: "shell", Command: "go test ./..."}, Config{
		ReadOnlyPrefixes: []string{"go test"},
	})
	if result.Verdict != VerdictSafe {
		t.Fatalf("verdict = %q, want %q; reasons: %v", result.Verdict, VerdictSafe, result.Reasons)
	}
}

func TestEvaluateWithConfigAddsRiskyPrefixesAndSubstrings(t *testing.T) {
	for _, command := range []string{"docker system prune", "echo token | pbcopy"} {
		t.Run(command, func(t *testing.T) {
			result := EvaluateWithConfig(Request{ToolName: "shell", Command: command}, Config{
				RiskyPrefixes:   []string{"docker system prune"},
				RiskySubstrings: []string{"| pbcopy"},
			})
			if result.Verdict != VerdictRisky {
				t.Fatalf("verdict = %q, want %q; reasons: %v", result.Verdict, VerdictRisky, result.Reasons)
			}
		})
	}
}

func TestEvaluateWithConfigKeepsRiskyBeforeReadOnly(t *testing.T) {
	result := EvaluateWithConfig(Request{ToolName: "shell", Command: "rm -rf /tmp/build"}, Config{
		ReadOnlyPrefixes: []string{"rm"},
	})
	if result.Verdict != VerdictRisky {
		t.Fatalf("verdict = %q, want %q", result.Verdict, VerdictRisky)
	}
}
