package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/romana/rlog"
)

func CleanUpExpiredPastes() {
	pastes, err := PasteDB.Query(
		"select ID, PasteName, Expire, Privacy, BurnAfter from pastes where Expire < datetime('now') or BurnAfter <= ReadCount and BurnAfter > 0",
	)
	if err != nil {
		rlog.Error("Could not clean up expired pastes", err)
	}
	defer pastes.Close()
	for pastes.Next() {
		var paste Paste
		err = pastes.Scan(
			&paste.ID,
			&paste.PasteName,
			&paste.Expire,
			&paste.Privacy,
			&paste.BurnAfter,
		)
		if err != nil {
			panic(err)
		}
		rlog.Info("Cleaning up expired paste", paste.PasteName)
	}
}
