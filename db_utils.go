package main

import (
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func IsSamePassword(pasteName, password string) bool {
	hashed := HashPassword(password)
	row := PasteDB.QueryRow(
		"select case when HashedPassword = ? then 1 else 0 end from pastes where paste_name = ?",
		hashed,
		pasteName,
	)
	var same int
	err := row.Scan(&same)
	if err != nil {
		panic(err)
	}
	return same == 1
}

func MakePastePointers(paste *Paste, scanVariables []string) []interface{} {
	pastePointers := make([]interface{}, len(scanVariables))
	val := reflect.ValueOf(paste).Elem()
	for i, variable := range scanVariables {
		pastePointers[i] = val.FieldByName(variable).Addr().Interface()
	}
	return pastePointers
}
