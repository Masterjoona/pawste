package database

import (
	"database/sql"
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
)

func UpdateReadCount(pasteName string) {
	amount := 1
	if !config.Vars.ReadCount {
		amount = 0
	}
	_, err := PasteDB.Exec(
		"update pastes set ReadCount = ReadCount + ?, ReadLast = strftime('%s', 'now') where PasteName = ?",
		amount,
		pasteName,
	)
	if err != nil {
		config.Logger.Error(err)
	}

	burnIfNeeded(pasteName)
}

func burnIfNeeded(pasteName string) {
	row := PasteDB.QueryRow(
		"select case when BurnAfter <= ReadCount and BurnAfter > 0 then 1 else 0 end from pastes where PasteName = ?",
		pasteName,
	)
	var burned int
	err := row.Scan(&burned)
	if err != nil {
		config.Logger.Error(err)
	}
	if burned == 1 {
		cleanUpExpiredPastes()
	}
}

func updatePasteContent(pasteName, content string) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			config.Logger.Error(err)
		}
		err = tx.Commit()
		if err != nil {
			config.Logger.Error(err)
		}
	}()
	stmt, err := tx.Prepare(`
		update pastes set
			Content = ?,
			UrlRedirect = ?,
			UpdatedAt = strftime('%s', 'now')
		where PasteName = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		content,
		utils.IsContentJustUrl(content),
		pasteName,
	)
	return err
}

func updatePasteFiles(pasteName string, newPaste paste.PasteUpdate) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			config.Logger.Error(err)
		}
		err = tx.Commit()
		if err != nil {
			config.Logger.Error(err)
		}
	}()

	for _, file := range newPaste.RemovedFiles {
		err = deleteFile(tx, pasteName, file)
		if err != nil {
			return err
		}
	}

	if len(newPaste.Files) > 0 {
		err = os.MkdirAll(config.Vars.DataDir+pasteName, 0755)
		if err != nil {
			return err
		}
	}

	for _, file := range newPaste.Files {
		err = insertFile(tx, pasteName, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteFile(tx *sql.Tx, pasteName, fileName string) error {
	stmt, err := tx.Prepare(`
		delete from files where PasteName = ? and Name = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		pasteName,
		fileName,
	)

	if err != nil {
		return err
	}

	return os.Remove(config.Vars.DataDir + pasteName + "/" + fileName)
}

func insertFile(tx *sql.Tx, pasteName string, file paste.File) error {
	stmt, err := tx.Prepare(`
		insert into files(PasteName, Name, Size, ContentType)
		values (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		pasteName,
		file.Name,
		file.Size,
		file.ContentType,
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		config.Vars.DataDir+pasteName+"/"+file.Name,
		file.Blob,
		0644,
	)
}

func UpdatePaste(pasteName string, newPaste paste.PasteUpdate) error {
	err := updatePasteContent(pasteName, newPaste.Content)
	if err != nil {
		return err
	}

	return updatePasteFiles(pasteName, newPaste)
}
