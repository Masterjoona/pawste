package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func RemovePaste(URLHref string) {
	if PasteDB == nil {
		panic("database connection is nil")
	}
	tx, err := PasteDB.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("delete from pastes where url = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(URLHref)
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
