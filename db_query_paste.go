package main

import (
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

func GetAllPastes() []Paste {
	CleanUpExpiredPastes()

	rows, err := PasteDB.Query("select paste_name, expire, privacy, burn_after from pastes")
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
	rows, err := PasteDB.Query(
		"select id, paste_name, expire, privacy, read_count, read_last, burn_after, content, url_redirect, syntax, hashed_password, created_at, updated_at FROM pastes WHERE paste_name = ?",
		pasteName,
	)
	if err != nil {
		return Paste{}, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(
			&paste.ID,
			&paste.PasteName,
			&paste.Expire,
			&paste.Privacy,
			&paste.ReadCount,
			&paste.ReadLast,
			&paste.BurnAfter,
			&paste.Content,
			&paste.UrlRedirect,
			&paste.Syntax,
			&paste.HashedPassword,
			&paste.CreatedAt,
			&paste.UpdatedAt,
		)
		if err != nil {
			return Paste{}, err
		}
		return paste, nil
	}
	return Paste{}, errors.New("paste not found")
}

func GetPublicPastes() []Paste {
	CleanUpExpiredPastes()
	rows, err := PasteDB.Query(
		"select paste_name, expire, read_count from pastes where privacy = 'public'",
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var paste Paste
		err = rows.Scan(&paste.PasteName, &paste.Expire, &paste.ReadCount)
		if err != nil {
			panic(err)
		}
		pastes = append(pastes, paste)
	}
	return pastes
}
