package mytracer

type Config struct {
	Enable         bool   `mapstructure:"enable"`
	ExporterTarget string `mapstructure:"exp_target"`
	// "host.docker.internal:4317"
}
