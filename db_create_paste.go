package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePaste(paste Paste) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}
	defer rollbackAndClose(tx, &err)

	stmt, err := tx.Prepare(`
		INSERT INTO pastes(PasteName, Expire, Privacy, ReadCount, ReadLast, BurnAfter, Content, UrlRedirect, Syntax, HashedPassword, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	encrypt := (paste.Privacy == "private" || paste.Privacy == "secret") && paste.HashedPassword != ""
	if encrypt {
		println("Encrypting paste content")
		encryptContent(&paste)
	}

	result, err := stmt.Exec(
		paste.PasteName,
		paste.Expire,
		paste.Privacy,
		paste.ReadCount,
		paste.ReadLast,
		paste.BurnAfter,
		paste.Content,
		paste.UrlRedirect,
		paste.Syntax,
		HashPassword(paste.HashedPassword), // finally hashed
		paste.CreatedAt,
		paste.UpdatedAt,
	)
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if len(paste.Files) > 0 {
		err = SaveFiles(tx, paste.Files, lastInsertID, paste.HashedPassword, encrypt)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveFiles(tx *sql.Tx, files []File, pasteID int64, password string, encrypt bool) error {
	for _, file := range files {
		if encrypt {
			encryptFile(&file, password)
		}

		_, err := tx.Exec(`
			INSERT INTO files(ID, Name, Size, Blob)
			VALUES (?, ?, ?, ?)
		`, pasteID, file.Name, file.Size, file.Blob)
		if err != nil {
			return err
		}
	}

	return nil
}

func rollbackAndClose(tx *sql.Tx, err *error) {
	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println("Failed to rollback transaction:", rollbackErr)
		}
		return
	}

	if commitErr := tx.Commit(); commitErr != nil {
		*err = commitErr
	}
}

func encryptContent(paste *Paste) {
	encryptedText, err := EncryptText(paste.Content, paste.HashedPassword)
	if err != nil {
		log.Println("Failed to encrypt content:", err)
		return
	}
	paste.Content = encryptedText
}

func encryptFile(file *File, password string) {
	err := Encrypt(file, password)
	if err != nil {
		log.Println("Failed to encrypt file:", err)
		return
	}

}
