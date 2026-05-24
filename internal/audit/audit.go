package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Path            string
	IncludeRawInput bool
}

type Writer struct {
	cfg Config
}

type Summary struct {
	ToolName string `json:"tool_name,omitempty"`
	Command  string `json:"command,omitempty"`
	CWD      string `json:"cwd,omitempty"`
}

type Record struct {
	Timestamp    time.Time `json:"timestamp"`
	Decision     string    `json:"decision"`
	Reason       string    `json:"reason,omitempty"`
	ProviderType string    `json:"provider_type,omitempty"`
	Model        string    `json:"model,omitempty"`
	Summary      Summary   `json:"summary"`
	RawInput     []byte    `json:"-"`
}

func New(cfg Config) Writer {
	return Writer{cfg: cfg}
}

func (w Writer) Write(record Record) error {
	if w.cfg.Path == "" {
		return nil
	}
	if record.Timestamp.IsZero() {
		record.Timestamp = time.Now().UTC()
	}
	line := auditLine{
		Timestamp:    record.Timestamp,
		Decision:     record.Decision,
		Reason:       record.Reason,
		ProviderType: record.ProviderType,
		Model:        record.Model,
		Summary:      record.Summary,
	}
	if len(record.RawInput) > 0 {
		sum := sha256.Sum256(record.RawInput)
		line.RawInputSHA256 = hex.EncodeToString(sum[:])
	}
	if w.cfg.IncludeRawInput {
		line.RawInput = string(record.RawInput)
	}

	encoded, err := json.Marshal(line)
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')

	file, err := os.OpenFile(w.cfg.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(encoded)
	return err
}

type auditLine struct {
	Timestamp      time.Time `json:"timestamp"`
	Decision       string    `json:"decision"`
	Reason         string    `json:"reason,omitempty"`
	ProviderType   string    `json:"provider_type,omitempty"`
	Model          string    `json:"model,omitempty"`
	Summary        Summary   `json:"summary"`
	RawInputSHA256 string    `json:"raw_input_sha256,omitempty"`
	RawInput       string    `json:"raw_input,omitempty"`
}
