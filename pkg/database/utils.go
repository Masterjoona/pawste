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

func pasteExists(name string) bool {
	var exists bool
	err := PasteDB.QueryRow("select exists(select 1 from pastes where PasteName = ?)", name).
		Scan(&exists)
	if err != nil {
		config.Logger.Error("Could not check if paste exists", err)
		return false
	}
	return exists
}

func HashPassword(password string) string {
	hash := pbkdf2.Key([]byte(password), []byte(config.Vars.Salt), 10000, 32, sha256.New)
	return fmt.Sprintf("%x", hash)
}
