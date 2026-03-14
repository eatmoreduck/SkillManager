package service

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"skillmanager/internal/model"
	"skillmanager/internal/repository"
)

// ConfigService manages persisted application settings and shared clients.
type ConfigService struct {
	configRepo repository.ConfigRepository
}

func NewConfigService(configRepo repository.ConfigRepository) *ConfigService {
	return &ConfigService{configRepo: configRepo}
}

func (s *ConfigService) GetConfig() (*model.Config, error) {
	return s.configRepo.LoadConfig()
}

func (s *ConfigService) UpdateProxy(proxy model.ProxyConfig) error {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return err
	}
	cfg.Proxy = proxy
	return s.configRepo.SaveConfig(cfg)
}

func (s *ConfigService) GetHTTPClient() *http.Client {
	cfg, err := s.configRepo.LoadConfig()
	if err != nil {
		return &http.Client{Timeout: 30 * time.Second}
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if cfg.Proxy.Enabled {
		proxyURL, err := url.Parse(cfg.Proxy.URL())
		if err != nil {
			return &http.Client{
				Timeout:   30 * time.Second,
				Transport: transport,
			}
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}
