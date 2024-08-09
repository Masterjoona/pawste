package database

import (
	"database/sql"
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

func CreatePaste(newPaste paste.Paste) error {
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

	encrypt := (newPaste.Privacy == "private" || newPaste.Privacy == "secret") &&
		newPaste.Password != ""
	if encrypt {
		newPaste.Content, err = paste.EncryptText(newPaste.Password, newPaste.Content)
		if err != nil {
			return err
		}
	}

	NewPassword := utils.Ternary((encrypt || newPaste.Privacy == "readonly"), HashPassword(newPaste.Password), "")

	_, err = stmt.Exec(
		newPaste.PasteName,
		newPaste.Expire,
		newPaste.Privacy,
		newPaste.NeedsAuth,
		newPaste.ReadCount,
		newPaste.ReadLast,
		newPaste.BurnAfter,
		newPaste.Content,
		newPaste.UrlRedirect,
		newPaste.Syntax,
		NewPassword,
		newPaste.CreatedAt,
		newPaste.UpdatedAt,
	)

	if err != nil {
		return err
	}

	err = saveFiles(tx, &newPaste, encrypt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func saveFiles(tx *sql.Tx, newPaste *paste.Paste, encrypt bool) error {
	for _, file := range newPaste.Files {
		if config.Vars.AnonymiseFileNames {
			file.Name = createShortFileName(newPaste.PasteName) // we dont have collision problems
		}
		if encrypt {
			err := paste.Encrypt(newPaste.Password, &file.Blob)
			if err != nil {
				config.Logger.Error("Failed to encrypt file:", err)
				return err
			}
		}
		err := saveFileToDisk(&file, newPaste.PasteName)
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
			newPaste.PasteName,
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
