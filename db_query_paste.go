package main

import (
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func getPastes(addQuery string) []Paste {
	CleanUpExpiredPastes()
	rows, err := PasteDB.Query(
		"select paste_name, expire, privacy, burn_after from pastes " + addQuery,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var paste Paste
		if err := rows.Scan(&paste.PasteName, &paste.Expire, &paste.Privacy, &paste.BurnAfter); err != nil {
			panic(err)
		}
		pastes = append(pastes, paste)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("Pastes: %v\n", pastes)
	return pastes
}

func GetAllPastes() []Paste {
	return getPastes(
		"",
	)
}

func GetPublicPastes() []Paste {
	return getPastes(
		"where privacy = 'public' and url_redirect = '0'",
	)
}

func GetPublicRedirects() []Paste {
	return getPastes(
		"where privacy = 'public' and url_redirect != '0'",
	)
}

// i do NOT like this, i wish i could assign the things i need to the paste struct

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
