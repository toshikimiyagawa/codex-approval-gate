package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriterDisabledDoesNothing(t *testing.T) {
	writer := New(Config{})

	if err := writer.Write(Record{Decision: "ask"}); err != nil {
		t.Fatal(err)
	}
}

func TestWriterWritesOneJSONLRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	writer := New(Config{Path: path})

	err := writer.Write(Record{
		Decision:     "allow",
		Reason:       "read-only",
		ProviderType: "openai",
		Model:        "local-model",
		Summary: Summary{
			ToolName: "shell",
			Command:  "git status",
			CWD:      "/tmp/project",
		},
		RawInput: []byte(`{"command":"git status"}`),
	})
	if err != nil {
		t.Fatal(err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var line map[string]any
	if err := json.Unmarshal(contents, &line); err != nil {
		t.Fatal(err)
	}
	if line["decision"] != "allow" {
		t.Fatalf("decision = %v, want allow", line["decision"])
	}
	if line["raw_input"] != nil {
		t.Fatalf("raw_input = %v, want omitted", line["raw_input"])
	}
	if line["raw_input_sha256"] == "" {
		t.Fatal("raw_input_sha256 was empty")
	}
}

func TestWriterIncludesRawInputWhenConfigured(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	writer := New(Config{Path: path, IncludeRawInput: true})

	if err := writer.Write(Record{Decision: "ask", RawInput: []byte(`{"x":1}`)}); err != nil {
		t.Fatal(err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var line map[string]any
	if err := json.Unmarshal(contents, &line); err != nil {
		t.Fatal(err)
	}
	if line["raw_input"] != `{"x":1}` {
		t.Fatalf("raw_input = %v, want raw json string", line["raw_input"])
	}
}
