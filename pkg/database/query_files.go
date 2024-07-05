package database

import (
	"errors"
	"strings"

	"github.com/Masterjoona/pawste/pkg/paste"
	_ "github.com/mattn/go-sqlite3"
)

func queryFiles(addQuery string, valueArgs []string, scanVariables []string) []paste.File {
	CleanUpExpiredPastes()

	valueInterfaces := make([]interface{}, len(valueArgs))
	for i, v := range valueArgs {
		valueInterfaces[i] = v
	}
	query := "select " + strings.Join(scanVariables, ", ") + " from files " + addQuery
	rows, err := PasteDB.Query(
		query,
		valueInterfaces...,
	)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var files []paste.File

	for rows.Next() {
		var file paste.File
		if err := rows.Scan(MakeFilePointers(&file, scanVariables)...); err != nil {
			panic(err)
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return files
}

func GetFiles(pasteName string) []paste.File {
	return queryFiles(
		"where PasteName = ?",
		[]string{pasteName},
		[]string{
			"ID",
			"Name",
			"Size",
		},
	)
}

func GetFile(pasteName string, fileName string) (paste.File, error) {
	files := queryFiles(
		"where PasteName = ? and Name = ?",
		[]string{pasteName, fileName},
		[]string{
			"ID",
			"Name",
			"Size",
		},
	)
	if len(files) == 0 {
		return paste.File{}, errors.New("file not found")
	}
	return files[0], nil
}
