package binding

import (
	"skillmanager/internal/model"
	"skillmanager/internal/service"
)

type ConfigBinding struct {
	service *service.ConfigService
}

type ProxyConfigResponse struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewConfigBinding(svc *service.ConfigService) *ConfigBinding {
	return &ConfigBinding{service: svc}
}

func (b *ConfigBinding) GetConfig() (*model.Config, error) {
	return b.service.GetConfig()
}

func (b *ConfigBinding) UpdateProxy(proxy model.ProxyConfig) error {
	return b.service.UpdateProxy(proxy)
}

func (b *ConfigBinding) GetProxy() (*ProxyConfigResponse, error) {
	config, err := b.service.GetConfig()
	if err != nil {
		return nil, err
	}

	return &ProxyConfigResponse{
		Enabled:  config.Proxy.Enabled,
		Type:     config.Proxy.Type,
		Host:     config.Proxy.Host,
		Port:     config.Proxy.Port,
		Username: config.Proxy.Username,
		Password: config.Proxy.Password,
	}, nil
}

func (b *ConfigBinding) SetProxy(enabled bool, proxyType string, host string, port int, username string, password string) error {
	return b.service.UpdateProxy(model.ProxyConfig{
		Enabled:  enabled,
		Type:     proxyType,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
}
