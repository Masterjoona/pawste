package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func CreatePaste(paste Paste, password string) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			panic(err)
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO pastes(paste_name, expire, privacy, read_count, read_last, burn_after, content, url_redirect, syntax, hashed_password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

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
		paste.HashedPassword,
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

	// Save files if any
	if len(paste.Files) > 0 {
		err = SaveFiles(paste, lastInsertID, password)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveFiles(paste Paste, insertId int64, password string) error {
	privacy := paste.Privacy
	encrypt := (privacy == "private" || privacy == "secret") && password != ""

	for _, file := range paste.Files {
		if encrypt {
			blob, err := Encrypt(file, password)
			if err != nil {
				return err
			}
			file.Blob = blob
		}
		_, err := PasteDB.Exec(`
			INSERT INTO files(paste_id, file_name, file_size, file_blob)
			VALUES (?, ?, ?, ?)
		`, insertId, file.Name, file.Size, file.Blob)
		if err != nil {
			return err
		}
	}

	return nil
}
