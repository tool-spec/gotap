package logging

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
)

const (
	EnvRunID     = "GOTAP_RUN_ID"
	EnvLogFile   = "GOTAP_LOG_FILE"
	EnvLogFormat = "GOTAP_LOG_FORMAT"
	EnvLogLevel  = "GOTAP_LOG_LEVEL"

	LogFormatJSONL = "jsonl"
	DefaultLevel   = "info"
)

type RunContext struct {
	RunID     string `json:"run_id"`
	LogFile   string `json:"log_file"`
	LogFormat string `json:"log_format"`
	LogLevel  string `json:"log_level"`
}

func NewRunContext(outputFolder string) (RunContext, error) {
	runID, err := randomRunID()
	if err != nil {
		return RunContext{}, err
	}

	return RunContext{
		RunID:     runID,
		LogFile:   filepath.Join(outputFolder, "_logs.jsonl"),
		LogFormat: LogFormatJSONL,
		LogLevel:  DefaultLevel,
	}, nil
}

func (c RunContext) Env() []string {
	return []string{
		EnvRunID + "=" + c.RunID,
		EnvLogFile + "=" + c.LogFile,
		EnvLogFormat + "=" + c.LogFormat,
		EnvLogLevel + "=" + c.LogLevel,
	}
}

func randomRunID() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
