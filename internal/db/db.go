package db

import (
	"database/sql"
	"time"
)

// New creates a new database connection pool.
func New(url string, maxOpenConn, maxIdleConn int, maxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxIdleTime(maxIdleTime)

	return db, nil
}
