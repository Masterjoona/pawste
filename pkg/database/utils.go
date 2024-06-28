package database

import (
	"crypto/sha256"
	"fmt"
	"reflect"

	"github.com/Masterjoona/pawste/pkg/paste"
	_ "github.com/mattn/go-sqlite3"
)

func MakePastePointers(paste *paste.Paste, scanVariables []string) []interface{} {
	pastePointers := make([]interface{}, len(scanVariables))
	val := reflect.ValueOf(paste).Elem()
	for i, variable := range scanVariables {
		pastePointers[i] = val.FieldByName(variable).Addr().Interface()
	}
	return pastePointers
}

func pasteExists(name string) bool {
	var exists bool
	err := PasteDB.QueryRow("select exists(select 1 from pastes where PasteName = ?)", name).
		Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write(paste.SecurePassword(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
