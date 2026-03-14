package model

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrAlreadyExists    = errors.New("resource already exists")
	ErrInvalidConfig    = errors.New("invalid config")
	ErrAgentNotEnabled  = errors.New("agent is not enabled")
	ErrSkillNotAssigned = errors.New("skill is not assigned to any enabled agent")
)

type ConfigLoadError struct {
	Path string
	Err  error
}

func (e *ConfigLoadError) Error() string {
	return "load config " + e.Path + ": " + e.Err.Error()
}

func (e *ConfigLoadError) Unwrap() error {
	return e.Err
}

type ConfigSaveError struct {
	Path string
	Err  error
}

func (e *ConfigSaveError) Error() string {
	return "save config " + e.Path + ": " + e.Err.Error()
}

func (e *ConfigSaveError) Unwrap() error {
	return e.Err
}
