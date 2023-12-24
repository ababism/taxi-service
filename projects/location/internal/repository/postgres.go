package repository

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	DBName     string `mapstructure:"dbname"`
	Password   string `mapstructure:"password"`
	SSLMode    string `mapstructure:"sslmode"`
	Migrations string `mapstructure:"migration"`
}

func NewPostgresDB(cfg *Config) (*sqlx.DB, error) {
	dbParams := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode,
	)
	db, err := sqlx.Open("postgres", dbParams)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("PostgresDB successfully up")

	return db, err
}
