package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

func GetAllPastes() []Paste {
	CleanUpExpiredPastes()

	rows, err := PasteDB.Query("select url, expire, privacy, burn_after from pastes")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var paste Paste
		err = rows.Scan(&paste.PasteName, &paste.Expire, &paste.Privacy, &paste.BurnAfter)
		if err != nil {
			panic(err)
		}
		pastes = append(pastes, paste)
	}
	return pastes
}

func GetPasteByName(pasteName string) (Paste, error) {
	CleanUpExpiredPastes()
	var paste Paste
	err := PasteDB.QueryRow("select id, url, expire, privacy, burn_after, content, syntax, hashed_password FROM pastes WHERE url = ?", pasteName).
		Scan(&paste.ID, &paste.PasteName, &paste.Expire, &paste.Privacy, &paste.BurnAfter, &paste.Content, &paste.Syntax, &paste.HashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return Paste{}, errors.New("paste not found")
		}
		return Paste{}, err
	}
	return paste, nil
}

func GetPublicPastes() []Paste {
	CleanUpExpiredPastes()
	rows, err := PasteDB.Query("select url, expire, burn_after from pastes where privacy = 'public'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var paste Paste
		err = rows.Scan(&paste.PasteName, &paste.Expire, &paste.BurnAfter)
		if err != nil {
			panic(err)
		}
		pastes = append(pastes, paste)
	}
	return pastes
}
