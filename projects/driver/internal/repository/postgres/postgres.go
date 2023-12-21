package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/config"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	dbParams := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.DBName,
		cfg.Postgres.Password,
		cfg.Postgres.SSLMode,
	)
	fmt.Println("dbParams", dbParams)
	db, err := sqlx.Open("postgres", dbParams)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("HI")
	return db, err
}
