package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CleanUpExpiredPastes() {
	pastes, err := PasteDB.Query(
		"select id, paste_name, expire, privacy, burn_after from pastes where expire < datetime('now') or burn_after <= read_count",
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
		println(fmt.Sprintf("Expired paste: %v", paste))
	}
}
