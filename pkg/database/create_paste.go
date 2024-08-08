package database

import (
	"database/sql"
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

func CreatePaste(paste paste.Paste) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}
	defer rollbackAndClose(tx, &err)

	stmt, err := tx.Prepare(`
		INSERT INTO pastes(PasteName, Expire, Privacy, NeedsAuth, ReadCount, ReadLast, BurnAfter, Content, UrlRedirect, Syntax, Password, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	encrypt := (paste.Privacy == "private" || paste.Privacy == "secret") &&
		paste.Password != ""
	if encrypt {
		err = paste.EncryptText(paste.Password)
		if err != nil {
			return err
		}
	}

	NewPassword := utils.Ternary((encrypt || paste.Privacy == "readonly"), HashPassword(paste.Password), "")

	_, err = stmt.Exec(
		paste.PasteName,
		paste.Expire,
		paste.Privacy,
		paste.NeedsAuth,
		paste.ReadCount,
		paste.ReadLast,
		paste.BurnAfter,
		paste.Content,
		paste.UrlRedirect,
		paste.Syntax,
		NewPassword,
		paste.CreatedAt,
		paste.UpdatedAt,
	)

	if err != nil {
		return err
	}

	err = saveFiles(tx, &paste, encrypt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func saveFiles(tx *sql.Tx, paste *paste.Paste, encrypt bool) error {
	for _, file := range paste.Files {
		if encrypt {
			err := file.Encrypt(paste.Password)
			if err != nil {
				config.Logger.Error("Failed to encrypt file:", err)
				return err
			}
		}
		err := saveFileToDisk(&file, paste.PasteName)
		if err != nil {
			config.Logger.Error("Failed to save file to disk:", err)
			return err
		}

		stmt, err := tx.Prepare(`
            INSERT INTO files(PasteName, Name, Size, ContentType)
            VALUES (?, ?, ?, ?)
        `)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			paste.PasteName,
			file.Name,
			file.Size,
			file.ContentType,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveFileToDisk(file *paste.File, pasteName string) error {
	err := os.MkdirAll(config.Vars.DataDir+pasteName, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(
		config.Vars.DataDir+pasteName+"/"+file.Name,
		file.Blob,
		0644,
	)
}

func rollbackAndClose(tx *sql.Tx, err *error) {
	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			config.Logger.Error("Failed to rollback transaction:", rollbackErr)
		}
		return
	}

	if commitErr := tx.Commit(); commitErr != nil {
		*err = commitErr
	}
}
