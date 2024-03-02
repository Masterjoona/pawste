package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func IsSamePassword(pasteName, password string) bool {
	hashed := HashPassword(password)
	row := PasteDB.QueryRow(
		"select case when hashed_password = ? then 1 else 0 end from pastes where paste_name = ?",
		hashed,
		pasteName,
	)
	var same int
	err := row.Scan(&same)
	if err != nil {
		panic(err)
	}
	return same == 1
}
