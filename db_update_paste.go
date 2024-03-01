package main

func UpdateReadCount(pasteName string) {
	_, err := PasteDB.Exec(
		"update pastes set read_count = read_count + 1, read_last = datetime('now') where paste_name = ?",
		pasteName,
	)
	if err != nil {
		panic(err)
	}
}

func IsAtBurnAfter(pasteName string) bool {
	row := PasteDB.QueryRow(
		"select burn_after, read_count from pastes where paste_name = ?",
		pasteName,
	)
	var burnAfter int
	var readCount int
	err := row.Scan(&burnAfter, &readCount)
	if err != nil {
		panic(err)
	}
	if burnAfter == 0 {
		return false
	}
	return readCount >= burnAfter
}
