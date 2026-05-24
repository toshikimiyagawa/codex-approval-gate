package codex

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	OutputModeCodex  = "codex"
	OutputModeSimple = "simple"
)

type PermissionRequest struct {
	ToolName string
	Command  string
	CWD      string
	Reason   string
	Raw      []byte
	Fields   map[string]any
}

func DecodePermissionRequest(input []byte) (PermissionRequest, error) {
	var fields map[string]any
	if err := json.Unmarshal(input, &fields); err != nil {
		return PermissionRequest{}, err
	}

	return PermissionRequest{
		ToolName: stringField(fields, "tool_name", "toolName", "tool"),
		Command:  stringField(fields, "command", "cmd"),
		CWD:      stringField(fields, "cwd", "current_working_directory"),
		Reason:   stringField(fields, "reason"),
		Raw:      append([]byte(nil), input...),
		Fields:   fields,
	}, nil
}

func EncodeResponse(decision string, outputMode string) ([]byte, error) {
	if !validDecision(decision) {
		return nil, fmt.Errorf("unknown decision %q", decision)
	}

	switch outputMode {
	case "", OutputModeCodex:
		return marshalLine(map[string]any{
			"hookSpecificOutput": map[string]string{
				"hookEventName":      "PermissionRequest",
				"permissionDecision": decision,
			},
		}, true)
	case OutputModeSimple:
		return marshalLine(map[string]string{"decision": decision}, false)
	default:
		return nil, errors.New("unknown codex output mode")
	}
}

func marshalLine(value any, indent bool) ([]byte, error) {
	var (
		output []byte
		err    error
	)
	if indent {
		output, err = json.MarshalIndent(value, "", "  ")
	} else {
		output, err = json.Marshal(value)
	}
	if err != nil {
		return nil, err
	}
	return append(output, '\n'), nil
}

func stringField(fields map[string]any, names ...string) string {
	for _, name := range names {
		value, ok := fields[name].(string)
		if ok {
			return value
		}
	}
	return ""
}

func validDecision(decision string) bool {
	return decision == "allow" || decision == "deny" || decision == "ask"
}
