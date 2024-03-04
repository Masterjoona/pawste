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
			(ID integer not null primary key, 
			PasteName text,
			Expire datetime, 
			Privacy text, 
			ReadCount integer, 
			ReadLast datetime,
			BurnAfter integer, 
			Content text,
			Syntax text,
			HashedPassword text,
			UrlRedirect integer,
			CreatedAt datetime, 
			UpdatedAt datetime
		);
		`
		_, err = sqldb.Exec(createPasteTable)
		if err != nil {
			log.Printf("%q: %s\n", err, createPasteTable)
			return nil
		}
		createFileTable := `
		create table files
			(ID integer not null primary key,
			PasteName text,
			Name text,
			Size integer,
			Blob blob
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
