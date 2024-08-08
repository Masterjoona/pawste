package database

import (
	"database/sql"
	"os"
	"time"

	"github.com/Masterjoona/pawste/pkg/config"
	_ "github.com/mattn/go-sqlite3"
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

	sqldb, err := sql.Open("sqlite3", config.Vars.DataDir+"pastes.db")
	if err != nil {
		config.Logger.Fatal("Could not open database", err)
	}
	if _, err := os.Stat(config.Vars.DataDir + "pastes.db"); os.IsNotExist(err) {
		_, err := sqldb.Exec(createPasteTable)
		if err != nil {
			config.Logger.Fatal("Could not create pastes table", err)
		}
		_, err = sqldb.Exec(createFileTable)
		if err != nil {
			config.Logger.Fatal("Could not create files table", err)
		}
		config.Logger.Info("Created new database")
	} else {
		config.Logger.Info("Loaded existing database")
	}
	PasteDB = sqldb

	go func() {
		cleanUpExpiredPastes()
		time.Sleep(1 * time.Hour)
	}()
}
