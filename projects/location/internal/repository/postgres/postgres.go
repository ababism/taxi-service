package postgres

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
	DBName     string `mapstructure:"db-name"`
	Password   string `mapstructure:"password"`
	SSLMode    string `mapstructure:"ssl-mode"`
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
	log.Println(fmt.Sprintf("New PostgresDB Params= %s migration_path=%s", dbParams, cfg.Migrations))
	db, err := sqlx.Open("postgres", dbParams)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("PostgresDB successfully up")

	//// Создаем объект миграции для Postgres
	//driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//m, err := migrate.NewWithDatabaseInstance(
	//	fmt.Sprintf("file:///%s", cfg.Migrations), // Путь к миграциям
	//	cfg.DBName, driver)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//m.Drop()
	//// Применяем миграции
	//if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	//	log.Fatal(err)
	//}
	//
	//log.Println("Migrations applied successfully")

	return db, err
}
