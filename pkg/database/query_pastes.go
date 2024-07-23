package database

import (
	"errors"
	"strings"

	"github.com/Masterjoona/pawste/pkg/paste"
	_ "github.com/mattn/go-sqlite3"
	"github.com/romana/rlog"
)

func queryPastes(addQuery string, valueArgs []string, scanVariables []string) []paste.Paste {
	cleanUpExpiredPastes()

	valueInterfaces := make([]interface{}, len(valueArgs))
	for i, v := range valueArgs {
		valueInterfaces[i] = v
	}
	query := "select " + strings.Join(scanVariables, ", ") + " from pastes " + addQuery
	rows, err := PasteDB.Query(
		query,
		valueInterfaces...,
	)

	if err != nil {
		rlog.Error("Could not query pastes", err)
		return nil
	}
	defer rows.Close()
	var pastes []paste.Paste

	for rows.Next() {
		var paste paste.Paste
		if err := rows.Scan(MakePastePointers(&paste, scanVariables)...); err != nil {
			rlog.Error("Could not scan paste", err)
			continue
		}
		pastes = append(pastes, paste)
	}
	if err := rows.Err(); err != nil {
		rlog.Error("Could not scan pastes", err)
		return nil
	}
	return pastes
}

func GetAllPastes() []paste.Paste {
	return queryPastes(
		"",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"ReadCount",
			"BurnAfter",
			"UrlRedirect",
			"ReadLast",
			"CreatedAt",
		})
}

func GetPublicPastes() []paste.Paste {
	return queryPastes(
		"where Privacy = 'public' and UrlRedirect = '0'",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"BurnAfter",
		},
	)
}

func GetPublicRedirects() []paste.Paste {
	return queryPastes(
		"where Privacy = 'public' and UrlRedirect != '0'",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"BurnAfter",
		},
	)
}

func GetAllPublicPastes() []paste.Paste {
	return queryPastes(
		"where Privacy = 'public'",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"ReadCount",
			"UrlRedirect",
		},
	)
}

func GetPasteByName(pasteName string) (paste.Paste, error) {
	pastes := queryPastes(
		"where PasteName = ?",
		[]string{pasteName},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"NeedsAuth",
			"ReadCount",
			"ReadLast",
			"BurnAfter",
			"Content",
			"UrlRedirect",
			"Syntax",
			"Password",
			"CreatedAt",
			"UpdatedAt",
		},
	)
	if len(pastes) == 0 {
		return paste.Paste{}, errors.New("paste not found")
	}
	return pastes[0], nil
}
