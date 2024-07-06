package database

import (
	"errors"
	"strings"

	"github.com/Masterjoona/pawste/pkg/paste"
	_ "github.com/mattn/go-sqlite3"
)

func queryPastes(addQuery string, valueArgs []string, scanVariables []string) []paste.Paste {
	CleanUpExpiredPastes()

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
		panic(err)
	}
	defer rows.Close()
	var pastes []paste.Paste

	for rows.Next() {
		var paste paste.Paste
		if err := rows.Scan(MakePastePointers(&paste, scanVariables)...); err != nil {
			panic(err)
		}
		pastes = append(pastes, paste)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return pastes
}

func GetAllPastes() paste.PasteLists {
	pastes := queryPastes(
		"where UrlRedirect = '0'",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"BurnAfter",
		})
	redirects := queryPastes(
		"where UrlRedirect != '0'",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"BurnAfter",
			"Syntax",
		})
	return paste.PasteLists{
		Pastes:    pastes,
		Redirects: redirects,
	}
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
			"ID",
			"PasteName",
			"Expire",
			"Privacy",
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
