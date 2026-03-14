package paths

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestPathsHelpers(t *testing.T) {
	t.Parallel()

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() error = %v", err)
	}
	cacheDir, err := GetCacheDir()
	if err != nil {
		t.Fatalf("GetCacheDir() error = %v", err)
	}
	skillsDir, err := GetSkillsDir()
	if err != nil {
		t.Fatalf("GetSkillsDir() error = %v", err)
	}
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}
	logPath, err := GetLogPath()
	if err != nil {
		t.Fatalf("GetLogPath() error = %v", err)
	}

	if !strings.HasSuffix(configDir, "skillmanager") {
		t.Fatalf("GetConfigDir() = %q, want suffix skillmanager", configDir)
	}
	if !strings.HasSuffix(cacheDir, "skillmanager") {
		t.Fatalf("GetCacheDir() = %q, want suffix skillmanager", cacheDir)
	}
	if !strings.HasSuffix(skillsDir, "skills") {
		t.Fatalf("GetSkillsDir() = %q, want suffix skills", skillsDir)
	}
	if !strings.HasSuffix(configPath, "config.yaml") {
		t.Fatalf("GetConfigPath() = %q, want suffix config.yaml", configPath)
	}
	if !strings.HasSuffix(logPath, "skillmanager.log") {
		t.Fatalf("GetLogPath() = %q, want suffix skillmanager.log", logPath)
	}

	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(strings.ToLower(configDir), "appdata") {
			t.Fatalf("windows config dir should contain AppData, got %q", configDir)
		}
	case "darwin":
		if !strings.Contains(configDir, "Library") {
			t.Fatalf("darwin config dir should contain Library, got %q", configDir)
		}
	case "linux":
		if !strings.Contains(configDir, ".config") {
			t.Fatalf("linux config dir should contain .config, got %q", configDir)
		}
	}
}

func TestEnsureDir(t *testing.T) {
	t.Parallel()

	target := t.TempDir() + "/nested/path"
	if err := EnsureDir(target); err != nil {
		t.Fatalf("EnsureDir() error = %v", err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("EnsureDir() did not create path: %v", err)
	}
}
