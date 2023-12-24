package http_server

import (
	"time"
)

const (
	defaultHost         = "localhost"
	defaultPort         = 8080
	defaultReadTimeout  = time.Second
	defaultWriteTimeout = time.Second
)

type Config struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// NewConfig use with dep on pkg/Config for configure rpc server
func NewConfig(prefix string) *Config {
	cfg := &Config{}
	// TODO Поставить флаги
	//config.StringVar(&cfg.Host, prefix+"host", defaultHost, "description")
	//config.IntVar(&cfg.Port, prefix+"port", defaultPort, "description")
	//config.DurationVar(&cfg.ReadTimeout, prefix+"read_timeout", defaultReadTimeout, "description")
	//config.DurationVar(&cfg.WriteTimeout, prefix+"write_timeout", defaultWriteTimeout, "description")

	return cfg
}

// NewDefaultConfig use when you need default server
func NewDefaultConfig() *Config {
	return &Config{
		Host:         defaultHost,
		Port:         defaultPort,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
}
