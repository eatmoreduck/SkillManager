package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCLIConfigCommand(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	t.Setenv("SKILLMANAGER_CONFIG", configPath)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := runCLI([]string{"config"}, &stdout, &stderr); err != nil {
		t.Fatalf("runCLI(config) error = %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Config path:") {
		t.Fatalf("runCLI(config) output = %q, want config path", output)
	}
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file was not created: %v", err)
	}
}
