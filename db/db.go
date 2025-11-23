package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	database := &DB{db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (db *DB) createTables() error {
	if err := db.createArtTable(); err != nil {
		return err
	}

	return nil
}
