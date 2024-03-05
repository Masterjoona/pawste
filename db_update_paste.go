package main

import (
	"errors"
)

func UpdateReadCount(pasteName string) {
	_, err := PasteDB.Exec(
		"update pastes set ReadCount = ReadCount + 1, ReadLast = datetime('now') where PasteName = ?",
		pasteName,
	)
	if err != nil {
		panic(err)
	}

	if isAtBurnAfter(pasteName) {
		println(pasteName, "should be deleted")
	}

}

func isAtBurnAfter(pasteName string) bool {
	row := PasteDB.QueryRow(
		"select case when BurnAfter <= ReadCount and BurnAfter > 0 then 1 else 0 end from pastes where PasteName = ?",
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
	// todo: check encryption level if password is needed
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
			Content = ?,
			Syntax = ?,
			UpdatedAt = datetime('now')
		where PasteName = ?
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
