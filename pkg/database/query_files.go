package database

import (
	"errors"
	"strings"

	"github.com/Masterjoona/pawste/pkg/paste"
	_ "github.com/mattn/go-sqlite3"
	"github.com/romana/rlog"
)

func queryFiles(addQuery string, valueArgs []string, scanVariables []string) []paste.File {
	cleanUpExpiredPastes()

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
		rlog.Error("Could not query files", err)
		return nil
	}
	defer rows.Close()
	var files []paste.File

	for rows.Next() {
		var file paste.File
		if err := rows.Scan(MakeFilePointers(&file, scanVariables)...); err != nil {
			rlog.Error("Could not scan file", err)
			continue
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		rlog.Error("Could not scan files", err)
		return nil
	}
	return files
}

func GetFiles(pasteName string) []paste.File {
	return queryFiles(
		"where PasteName = ?",
		[]string{pasteName},
		[]string{
			"Name",
			"Size",
			"ContentType",
		},
	)
}

func GetFile(pasteName string, fileName string) (paste.File, error) {
	files := queryFiles(
		"where PasteName = ? and Name = ?",
		[]string{pasteName, fileName},
		[]string{
			"Name",
			"Size",
			"ContentType",
		},
	)
	if len(files) == 0 {
		return paste.File{}, errors.New("file not found")
	}
	return files[0], nil
}
