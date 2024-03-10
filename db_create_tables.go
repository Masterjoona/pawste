package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/romana/rlog"
)

func CreateOrLoadDatabase(deleteOld bool) *sql.DB {
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

	createFileTable := `
	create table files
		(ID integer not null primary key,
		PasteName text,
		Name text,
		Size integer,
	);
	`
	if deleteOld {
		os.Remove("./pastes.db")
	}
	newDb := false
	sqldb, err := sql.Open("sqlite3", "./pastes.db")
	if err != nil {
		rlog.Critical("Could not open database", err)
	}
	if deleteOld || newDb {
		_, err = sqldb.Exec(createPasteTable)
		if err != nil {
			rlog.Critical("Could not create pastes table", err)
		}
		_, err = sqldb.Exec(createFileTable)
		if err != nil {
			rlog.Critical("Could not create files table", err)
		}
		rlog.Info("Created new database")
	} else {
		rlog.Info("Loaded existing database")
	}
	return sqldb
}
