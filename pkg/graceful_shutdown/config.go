package graceful_shutdown

import (
	"time"
)

const (
	defaultDelay           = 5 * time.Second
	defaultWaitTimeout     = 10 * time.Second
	defaultCallbackTimeout = 2 * time.Second
)

type Config struct {
	Delay           time.Duration `yaml:"delay"`
	WaitTimeout     time.Duration `yaml:"wait_timeout"`
	CallbackTimeout time.Duration `yaml:"callback_timeout"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Delay:           defaultDelay,
		WaitTimeout:     defaultWaitTimeout,
		CallbackTimeout: defaultCallbackTimeout,
	}
}
