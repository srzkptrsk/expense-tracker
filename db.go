// db
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	User     string
	Password string
	Schema   string
	Host     string
	Port     string
}

type Env struct {
	DB *sql.DB
}

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
