package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func CreatePaste(paste Paste) {
	if PasteDB == nil {
		panic("database connection is nil")
	}
	tx, err := PasteDB.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("insert into pastes(url, expire, privacy, burn_after, content, syntax, hashed_password) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	println(paste.Expire)
	_, err = stmt.Exec(paste.PasteName, paste.Expire, paste.Privacy, paste.BurnAfter, paste.Content, paste.Syntax, paste.HashedPassword)
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
