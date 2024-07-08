package database

import (
	"database/sql"

	"github.com/Masterjoona/pawste/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

var PasteDB *sql.DB

func init() {
	PasteDB = CreateOrLoadDatabase(config.Config.IUnderstandTheRisks)
}
