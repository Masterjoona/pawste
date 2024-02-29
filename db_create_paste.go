package main

import (
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePaste(paste Paste, password string) {
	tx, err := PasteDB.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("insert into pastes(paste_name, expire, privacy, read_count, read_last, burn_after, content, url_redirect, syntax, hashed_password, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(paste.PasteName, paste.Expire, paste.Privacy, paste.ReadCount, paste.ReadLast, paste.BurnAfter, paste.Content, paste.UrlRedirect, paste.Syntax, paste.HashedPassword, paste.CreatedAt, paste.UpdatedAt)
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	encrypt := (paste.Privacy == "private" || paste.Privacy == "secret") && password != ""
	if encrypt {
		println("Password:", password)
		//SaveFilesIfExists(paste.Files, paste.PasteName, encrypt, password)
		return
	}

	SaveFilesIfExists(paste.Files, paste.PasteName, encrypt, password)

}

func SaveFilesIfExists(files []File, pasteName string, encrypt bool, password string) {
	if len(files) == 0 || files[0].FileName == "" {
		return
	}
	_ = os.MkdirAll(Config.DataDir+pasteName, os.ModePerm)
	if encrypt {
		for index, file := range files {
			Encrypt(file.FileBlob, password)
		}
		return
	}
	for _, file := range files {
		println("Saving file:", file.Filename)
		src, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer src.Close()

		dst, err := os.Create(Config.DataDir + pasteName + "/" + file.Filename)
		if err != nil {
			panic(err)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			panic(err)
		}
	}
}
