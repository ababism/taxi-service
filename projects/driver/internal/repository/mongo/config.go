package mongo

type Config struct {
	Database string `mapstructure:"database" env:"DRIVER_MONGO_DATABASE"`
	Uri      string `mapstructure:"uri" env:"DRIVER_MONGO_URI"`
}

type ConfigMigrations struct {
	URI     string `mapstructure:"uri" env:"DRIVER_MONGO_MIGRATION_URI"`
	Path    string `mapstructure:"path" env:"DRIVER_MIGRATIONS_PATH"`
	Enabled bool   `mapstructure:"enabled" env:"DRIVER_MIGRATIONS_ENABLED"`
}
