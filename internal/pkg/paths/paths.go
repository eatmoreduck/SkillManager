package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const appName = "skillmanager"

func GetConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appName), nil
}

func GetCacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appName), nil
}

func GetSkillsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "skills"), nil
}

func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

func GetLogPath() (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "logs", "skillmanager.log"), nil
}

func Expand(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if strings.HasPrefix(path, "~") {
		return homedir.Expand(path)
	}

	return filepath.Clean(path), nil
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func HomePath(parts ...string) string {
	items := append([]string{"~"}, parts...)
	return filepath.Join(items...)
}

func PlatformName() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
