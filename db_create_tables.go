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
		createPasteTable := `
		create table pastes 
			(id integer not null primary key, 
			paste_name text,
			expire datetime, 
			privacy text, 
			read_count integer, 
			read_last datetime,
			burn_after integer, 
			content text,
			syntax text,
			hashed_password text,
			url_redirect integer,
			created_at datetime, 
			updated_at datetime
		);
		`
		_, err = sqldb.Exec(createPasteTable)
		if err != nil {
			log.Printf("%q: %s\n", err, createPasteTable)
			return nil
		}
		createFileTable := `
		create table files
			(id integer not null primary key,
			paste_name text,
			file_name text,
			file_size integer,
			file blob
		);
		`
		_, err = sqldb.Exec(createFileTable)
		if err != nil {
			log.Printf("%q: %s\n", err, createFileTable)
			return nil
		}
	}
	return sqldb
}
