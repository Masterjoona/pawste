package main

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
		"select if(burn_after <= read_count, 1, 0) from pastes where paste_name = ?",
		pasteName,
	)
	var burned int
	err := row.Scan(&burned)
	if err != nil {
		panic(err)
	}
	return burned == 1

}
