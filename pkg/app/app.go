package app

import (
	"fmt"
)

type Info struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"env"`
	Version     string `yaml:"version"`
}

const (
	production = "production"
	staging    = "staging"
	local      = "local"
)

func InitInfo(buildVersion string) *Info {
	info := &Info{}
	info.Version = buildVersion
	// TODO Ставим флаги
	//config.StringVar(&info.Name, "app.name", "unknown app", "description")
	//config.StringVar(&info.Environment, "app.env", "local", "description")
	//config.StringVar(&info.Owner, "app.owner", "unknown", "description")
	//config.StringVar(&info.Process, "app.process", "*", "comma separated processes to run. http/rpc/*...")

	return info
}

func (i *Info) Release() string {
	return fmt.Sprintf("%s-%s", i.Environment, i.Version)
}

// IsProduction defines is current app.env a "production"
func (i *Info) IsProduction() bool {
	return i.Environment == production
}

// IsStaging defines is current app.env a "staging"
func (i *Info) IsStaging() bool {
	return i.Environment == staging
}

// IsLocal defines is current app.env a "local"
func (i *Info) IsLocal() bool {
	return i.Environment == local
}
