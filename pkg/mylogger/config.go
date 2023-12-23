package mylogger

type Config struct {
	LevelLogger      string   `mapstructure:"level_logger"`
	LevelSentry      string   `mapstructure:"level_sentry"`
	Env              string   `mapstructure:"env"`
	OutputPaths      []string `mapstructure:"outputs"`
	ErrorOutputPaths []string `mapstructure:"error_outputs"`
	Encoding         string   `mapstructure:"encoding"`
	SentryDSN        string   `mapstructure:"sentry_dsn"`
}
