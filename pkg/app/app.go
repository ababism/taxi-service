package app

import (
	"fmt"
)

type Config struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"env"`
	Version     string `mapstructure:"version"`
}

const (
	production = "production"
	staging    = "staging"
	local      = "local"
)

func InitInfo(buildVersion string) *Config {
	info := &Config{}
	info.Version = buildVersion
	// TODO Ставим флаги
	//config.StringVar(&info.Name, "app.name", "unknown app", "description")
	//config.StringVar(&info.Environment, ".env", "local", "description")
	//config.StringVar(&info.Owner, "app.owner", "unknown", "description")
	//config.StringVar(&info.Process, "app.process", "*", "comma separated processes to run. http/rpc/*...")

	return info
}

func (i *Config) Release() string {
	return fmt.Sprintf("%s-%s", i.Environment, i.Version)
}

// IsProduction defines is current .env a "production"
func (i *Config) IsProduction() bool {
	return i.Environment == production
}

// IsStaging defines is current .env a "staging"
func (i *Config) IsStaging() bool {
	return i.Environment == staging
}

// IsLocal defines is current .env a "local"
func (i *Config) IsLocal() bool {
	return i.Environment == local
}
