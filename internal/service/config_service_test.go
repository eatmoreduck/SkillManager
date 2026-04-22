package service

import (
	"net/http"
	"testing"
	"time"

	"skillmanager/internal/model"
)

type mockConfigRepository struct {
	config *model.Config
	err    error
}

func (m *mockConfigRepository) LoadConfig() (*model.Config, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.config == nil {
		return &model.Config{Version: "1.0"}, nil
	}
	return m.config, nil
}

func (m *mockConfigRepository) SaveConfig(config *model.Config) error {
	m.config = config
	return m.err
}

func (m *mockConfigRepository) GetConfigPath() string {
	return "/tmp/config.yaml"
}

func TestConfigServiceGetHTTPClient(t *testing.T) {
	t.Parallel()

	svc := NewConfigService(&mockConfigRepository{
		config: &model.Config{
			Version: "1.0",
			Proxy: model.ProxyConfig{
				Enabled: true,
				Type:    "http",
				Host:    "127.0.0.1",
				Port:    8080,
			},
		},
	})

	client := svc.GetHTTPClient()
	if client == nil {
		t.Fatal("GetHTTPClient() = nil")
	}
	if client.Timeout != 30*time.Second {
		t.Fatalf("GetHTTPClient() timeout = %v, want %v", client.Timeout, 30*time.Second)
	}
	if _, ok := client.Transport.(*http.Transport); !ok {
		t.Fatalf("GetHTTPClient() transport = %T, want *http.Transport", client.Transport)
	}
}

func TestConfigServiceUpdateLanguage(t *testing.T) {
	t.Parallel()

	repo := &mockConfigRepository{
		config: &model.Config{
			Version:  "1.0",
			Language: "zh-CN",
		},
	}
	svc := NewConfigService(repo)

	if err := svc.UpdateLanguage("en-US"); err != nil {
		t.Fatalf("UpdateLanguage() error = %v", err)
	}

	if repo.config.Language != "en-US" {
		t.Fatalf("UpdateLanguage() language = %q, want %q", repo.config.Language, "en-US")
	}
}
