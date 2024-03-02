package main

import (
	"errors"
)

func UpdateReadCount(pasteName string) {
	_, err := PasteDB.Exec(
		"update pastes set read_count = read_count + 1, read_last = datetime('now') where paste_name = ?",
		pasteName,
	)
	if err != nil {
		panic(err)
	}

	if isAtBurnAfter(pasteName) {
		RemovePaste(pasteName)
	}

}

func isAtBurnAfter(pasteName string) bool {
	row := PasteDB.QueryRow(
		"select case when burn_after <= read_count and burn_after > 0 then 1 else 0 end from pastes where paste_name = ?",
		pasteName,
	)
	var burned int
	err := row.Scan(&burned)
	if err != nil {
		panic(err)
	}
	return burned == 1
}

func UpdatePaste(paste Paste, password string) error {
	if !IsSamePassword(paste.PasteName, password) {
		return errors.New("wrong password")
	}

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
			content = ?,
			syntax = ?,
			updated_at = datetime('now')
		where paste_name = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		paste.Content,
		paste.Syntax,
		paste.PasteName,
	)
	if err != nil {
		return err
	}

	// todo update files
	return nil
}
