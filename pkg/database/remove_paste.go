package database

import (
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

func removePaste(pasteName string) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("delete from pastes where PasteName = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(pasteName)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	//rlog.Info("Removed paste from database", pasteName)
	return nil
}

func removeFiles(pasteName string) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("delete from files where PasteName = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(pasteName)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	//rlog.Info("Removed files from db for paste", pasteName)

	return os.RemoveAll(config.Config.DataDir + pasteName)
}

func DeletePaste(pasteName string) error {
	err := removeFiles(pasteName)
	if err != nil {
		return err
	}

	return removePaste(pasteName)
}
