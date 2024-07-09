package handling

import "github.com/Masterjoona/pawste/pkg/database"

func isValidPassword(inputPassword, storedPassword string) bool {
	return database.HashPassword(inputPassword) == storedPassword
}
