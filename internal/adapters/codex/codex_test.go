package codex

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDecodePermissionRequestExtractsSummaryAndRawInput(t *testing.T) {
	input, err := os.ReadFile("testdata/permission_request.json")
	if err != nil {
		t.Fatal(err)
	}

	req, err := DecodePermissionRequest(input)
	if err != nil {
		t.Fatal(err)
	}

	if req.ToolName != "shell" {
		t.Fatalf("ToolName = %q, want shell", req.ToolName)
	}
	if req.Command != "git status --short" {
		t.Fatalf("Command = %q, want git status --short", req.Command)
	}
	if req.CWD != "/tmp/project" {
		t.Fatalf("CWD = %q, want /tmp/project", req.CWD)
	}
	if string(req.Raw) != string(input) {
		t.Fatal("raw input was not preserved")
	}
}

func TestEncodeResponseDefaultsToCodexHookFormat(t *testing.T) {
	got, err := EncodeResponse("allow", OutputModeCodex)
	if err != nil {
		t.Fatal(err)
	}

	var decoded struct {
		HookSpecificOutput struct {
			HookEventName      string `json:"hookEventName"`
			PermissionDecision string `json:"permissionDecision"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.HookSpecificOutput.HookEventName != "PermissionRequest" {
		t.Fatalf("hookEventName = %q, want PermissionRequest", decoded.HookSpecificOutput.HookEventName)
	}
	if decoded.HookSpecificOutput.PermissionDecision != "allow" {
		t.Fatalf("permissionDecision = %q, want allow", decoded.HookSpecificOutput.PermissionDecision)
	}
}

func TestEncodeResponseSupportsSimpleFormat(t *testing.T) {
	got, err := EncodeResponse("ask", OutputModeSimple)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"decision":"ask"}` + "\n"
	if string(got) != want {
		t.Fatalf("simple response = %s, want %s", got, want)
	}
}

func TestEncodeResponseRejectsUnknownDecision(t *testing.T) {
	if _, err := EncodeResponse("maybe", OutputModeCodex); err == nil {
		t.Fatal("expected error for unknown decision")
	}
}
