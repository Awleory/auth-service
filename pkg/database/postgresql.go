package database

import (
	"database/sql"
	"fmt"
)

type Config struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func Connection(cfg Config) (*sql.DB, error) {
	constring := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.SSLMode, cfg.Password)

	db, err := sql.Open("postgres", constring)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
