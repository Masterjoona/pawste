package database

import (
	"github.com/Masterjoona/pawste/pkg/paste"
)

func UpdateReadCount(pasteName string) {
	_, err := PasteDB.Exec(
		"update pastes set ReadCount = ReadCount + 1, ReadLast = datetime('now') where PasteName = ?",
		pasteName,
	)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	if burned == 1 {
		cleanUpExpiredPastes()
	}
}

func updatePasteContent(paste paste.Paste) error {
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
		paste.Content,
		paste.PasteName,
	)
	if err != nil {
		return err
	}
	return nil
}

func updatePasteFiles(paste paste.Paste) error {
	return nil
}

func UpdatePaste(paste paste.Paste) error {
	err := updatePasteContent(paste)
	if err != nil {
		return err
	}

	return updatePasteFiles(paste)
}
