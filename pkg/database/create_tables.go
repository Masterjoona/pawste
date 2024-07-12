package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/romana/rlog"
)

var PasteDB *sql.DB

func CreateOrLoadDatabase() {
	createPasteTable := `
	create table pastes 
		(ID integer not null primary key, 
		PasteName text,
		Expire integer, 
		Privacy text, 
		NeedsAuth integer,
		ReadCount integer, 
		ReadLast integer,
		BurnAfter integer, 
		Content text,
		Syntax text,
		Password text,
		UrlRedirect integer,
		CreatedAt integer, 
		UpdatedAt integer
	);
	`

	createFileTable := `
	create table files
		(ID integer not null primary key,
		PasteName text,
		Name text,
		Size integer,
		ContentType text
	);
	`

	sqldb, err := sql.Open("sqlite3", "./pastes.db")
	if err != nil {
		rlog.Critical("Could not open database", err)
	}
	if _, err := os.Stat("./pastes.db"); os.IsNotExist(err) {
		_, err := sqldb.Exec(createPasteTable)
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
	PasteDB = sqldb
}
