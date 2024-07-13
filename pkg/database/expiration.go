package database

import (
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/romana/rlog"
)

func cleanUpExpiredPastes() {
	pastes, err := PasteDB.Query(
		"select PasteName from pastes where Expire < strftime('%s', 'now') or BurnAfter <= ReadCount and BurnAfter > 0",
	)
	if err != nil {
		rlog.Error("Could not clean up expired pastes", err)
		return
	}
	defer pastes.Close()

	tx, err := PasteDB.Begin()
	if err != nil {
		rlog.Error("Could not start transaction", err)
		return
	}

	for pastes.Next() {
		var pasteName string
		err = pastes.Scan(&pasteName)
		if err != nil {
			rlog.Error("Could not scan pasteName", err)
			tx.Rollback()
			return
		}

		rlog.Info("Cleaning up paste with name " + pasteName)
		_, err = tx.Exec("delete from pastes where PasteName = ?", pasteName)
		if err != nil {
			rlog.Error("Could not delete paste with name "+pasteName, err)
			tx.Rollback()
			return
		}
		_, err = tx.Exec("delete from files where PasteName = ?", pasteName)
		if err != nil {
			rlog.Error("Could not delete files in db for paste with name "+pasteName, err)
			tx.Rollback()
			return
		}
		err = os.RemoveAll(config.Vars.DataDir + "/" + pasteName)
		if err != nil {
			rlog.Error("Could not delete files for paste with name "+pasteName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		rlog.Error("Could not commit transaction", err)
		return
	}

	rlog.Info("Expired pastes cleaned up successfully")
}
