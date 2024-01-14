package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// PgConfig
// структура с данными для подключения
type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

//NewPgConfig
//func NewPgConfig(username string, password string, host string, port string, database string, SSLMode string) *Config {
//	return &Config{Username: username, Password: password, Host: host, Port: port, Database: database, SSLMode: SSLMode}
//}

// NewClient
// инициализируем базу данных на основе данных из структуры
func NewPostgresDB(cfg *Config) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
