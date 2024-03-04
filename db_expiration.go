package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CleanUpExpiredPastes() {
	pastes, err := PasteDB.Query(
		"select ID, PasteName, Expire, Privacy, BurnAfter from pastes where Expire < datetime('now') or BurnAfter <= ReadCount and BurnAfter > 0",
	)
	if err != nil {
		panic(err)
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
		log.Printf("Cleaning up paste %s", paste.PasteName)
	}
}
