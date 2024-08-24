package database

import (
	"crypto/sha256"
	"fmt"
	"reflect"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
)

func MakePastePointers(paste *paste.Paste, scanVariables []string) []interface{} {
	pastePointers := make([]interface{}, len(scanVariables))
	val := reflect.ValueOf(paste).Elem()
	for i, variable := range scanVariables {
		pastePointers[i] = val.FieldByName(variable).Addr().Interface()
	}
	return pastePointers
}

func MakeFilePointers(paste *paste.File, scanVariables []string) []interface{} {
	filePointers := make([]interface{}, len(scanVariables))
	val := reflect.ValueOf(paste).Elem()
	for i, variable := range scanVariables {
		filePointers[i] = val.FieldByName(variable).Addr().Interface()
	}
	return filePointers
}

func exists(query string, args ...interface{}) bool {
	var exists bool
	err := PasteDB.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		config.Logger.Error("Could not check if exists", err)
		return false
	}
	return exists
}

func pasteExists(name string) bool {
	return exists("select exists(select 1 from pastes where PasteName = ?)", name)
}

func fileExists(pasteName, name string) bool {
	return exists("select exists(select 1 from files where PasteName = ? and Name = ?)", pasteName, name)
}

func HashPassword(password string) string {
	hash := pbkdf2.Key([]byte(password), []byte(config.Vars.Salt), 10000, 32, sha256.New)
	return fmt.Sprintf("%x", hash)
}
