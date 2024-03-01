package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CleanUpExpiredPastes() {
	/*_, err := PasteDB.Exec("delete from pastes where expire < datetime('now')")
	if err != nil {
		panic(err)
	}*/
	// for now we will just print the expired pastes
	pastes, err := PasteDB.Query("select * from pastes where expire < datetime('now')")
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
		println(fmt.Sprintf("Expired paste: %v", paste))
	}
}
