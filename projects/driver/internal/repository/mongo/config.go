package mongo

type MongoCfg struct {
	Database string `yaml:"database"          env:"DRIVER_MONGO_DATABASE" default:"driver"`
	Uri      string `yaml:"uri"               env:"DRIVER_MONGO_URI"`
}

type MigrationsCfg struct {
	URI     string `yaml:"uri"     env:"DRIVER_MONGO_MIGRATION_URI"`
	Path    string `yaml:"path"    env:"DRIVER_MIGRATIONS_PATH"`
	Enabled bool   `yaml:"enabled" env:"DRIVER_MIGRATIONS_ENABLED"   default:"false"`
}
