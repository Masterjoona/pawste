package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateOrLoadDatabase(deleteOld bool) *sql.DB {
	if deleteOld {
		os.Remove("./pastes.db")
	}

	sqldb, err := sql.Open("sqlite3", "./pastes.db")
	if err != nil {
		log.Fatal(err)
	}
	if deleteOld {
		sqlStmt := `
		create table pastes (id integer not null primary key, url text, expire datetime, privacy text, burn_after integer, content text, syntax text, hashed_password text);
		`
		_, err = sqldb.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return nil
		}
	}
	return sqldb
}
