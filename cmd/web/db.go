package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func openDB(driverName, dsn string) (*sql.DB, error) {
    db, err := sql.Open(driverName, dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}