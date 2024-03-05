package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/romana/rlog"
)

func RemovePaste(pasteName string) {
	tx, err := PasteDB.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("delete from pastes where PasteName = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(pasteName)
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	rlog.Info("Removed paste from database", pasteName)
}

func RemoveFiles(pasteName string) {
	tx, err := PasteDB.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("delete from files where PasteName = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(pasteName)
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	rlog.Info("Removed files from db for paste", pasteName)
}
