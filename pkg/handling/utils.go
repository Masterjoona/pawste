package handling

import "github.com/Masterjoona/pawste/pkg/database"

func verifyAccess(privacy bool, reqPasswd, hashedPasswd string) bool {
	return privacy && database.HashPassword(reqPasswd) != hashedPasswd
}
