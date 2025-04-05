package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func Connect(user, password, db, host string) *DB {
	con, err := sql.Open("postgres", "host=" + host + " user="+user+" password="+password+" dbname="+db+" port=5432 sslmode=disable")
	if err != nil {
		fmt.Errorf("Error while connecting to DB: ", err.Error())
	}

	if err := con.Ping(); err != nil {
		fmt.Errorf("Error while pining DB: ", err.Error())
	}

	return &DB{
		DB: con,
	}
}
