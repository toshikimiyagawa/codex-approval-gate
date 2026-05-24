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

	tool := mapField(fields, "tool")
	toolInput := mapField(tool, "input", "arguments", "args")
	if len(toolInput) == 0 {
		toolInput = mapField(fields, "tool_input", "toolInput", "input", "arguments", "args")
	}

	return PermissionRequest{
		ToolName: firstString(
			stringField(fields, "tool_name", "toolName"),
			stringField(tool, "name", "tool_name", "toolName"),
			stringField(fields, "tool"),
		),
		Command: firstString(
			stringField(fields, "command", "cmd"),
			stringField(toolInput, "command", "cmd"),
		),
		CWD: firstString(
			stringField(fields, "cwd", "current_working_directory"),
			stringField(toolInput, "cwd", "current_working_directory"),
		),
		Reason: stringField(fields, "reason"),
		Raw:    append([]byte(nil), input...),
		Fields: fields,
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

func mapField(fields map[string]any, names ...string) map[string]any {
	for _, name := range names {
		value, ok := fields[name].(map[string]any)
		if ok {
			return value
		}
	}
	return nil
}

func firstString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func validDecision(decision string) bool {
	return decision == "allow" || decision == "deny" || decision == "ask"
}
