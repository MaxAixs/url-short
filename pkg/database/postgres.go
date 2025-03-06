package database

import (
	"database/sql"
	"fmt"
	"log"
	"url-shortener/app/config"
)

func NewDB(cfg *config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Printf("connected to database: on host:%v port: %v\n", cfg.Host, cfg.Port)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}
