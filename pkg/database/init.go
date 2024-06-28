package database

import (
	"database/sql"

	"github.com/Masterjoona/pawste/pkg/shared/config"
	_ "github.com/mattn/go-sqlite3"
)

var PasteDB *sql.DB

func init() {
	PasteDB = CreateOrLoadDatabase(config.Config.IUnderstandTheRisks)
}
