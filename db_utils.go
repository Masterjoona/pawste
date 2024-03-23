package main

import (
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func MakePastePointers(paste *Paste, scanVariables []string) []interface{} {
	pastePointers := make([]interface{}, len(scanVariables))
	val := reflect.ValueOf(paste).Elem()
	for i, variable := range scanVariables {
		pastePointers[i] = val.FieldByName(variable).Addr().Interface()
	}
	return pastePointers
}
