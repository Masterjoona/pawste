package database

import (
	"database/sql"
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
	"github.com/romana/rlog"
)

func UpdateReadCount(pasteName string) {
	/*_, err := PasteDB.Exec(
		"update pastes set ReadCount = ReadCount + 1, ReadLast = datetime('now') where PasteName = ?",
		pasteName,
	)
	if err != nil {
		rlog.Error(err)
	}

	burnIfNeeded(pasteName)*/
	return
}

func burnIfNeeded(pasteName string) {
	row := PasteDB.QueryRow(
		"select case when BurnAfter <= ReadCount and BurnAfter > 0 then 1 else 0 end from pastes where PasteName = ?",
		pasteName,
	)
	var burned int
	err := row.Scan(&burned)
	if err != nil {
		rlog.Error(err)
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
			rlog.Error(err)
		}
		err = tx.Commit()
		if err != nil {
			rlog.Error(err)
		}
	}()

	stmt, err := tx.Prepare(`
		update pastes set
			Content = ?,
			UpdatedAt = datetime('now')
		where PasteName = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		content,
		pasteName,
	)
	return err
}

func updatePasteFiles(pasteName string, newPaste utils.PasteUpdate) error {
	tx, err := PasteDB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			rlog.Error(err)
		}
		err = tx.Commit()
		if err != nil {
			rlog.Error(err)
		}
	}()

	for _, file := range newPaste.RemovedFiles {
		err = deleteFile(tx, pasteName, file)
		if err != nil {
			return err
		}
	}

	if len(newPaste.Files) > 0 {
		err = os.MkdirAll(config.Config.DataDir+pasteName, 0755)
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

	return os.Remove(config.Config.DataDir + pasteName + "/" + fileName)
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
		config.Config.DataDir+pasteName+"/"+file.Name,
		file.Blob,
		0644,
	)
}

func UpdatePaste(pasteName string, newPaste utils.PasteUpdate) error {
	err := updatePasteContent(pasteName, newPaste.Content)
	if err != nil {
		return err
	}

	return updatePasteFiles(pasteName, newPaste)
}
