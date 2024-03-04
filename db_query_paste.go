package main

import (
	"errors"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func getPastes(addQuery string, valueArgs []string, scanVariables []string) []Paste {
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
	var pastes []Paste

	for rows.Next() {
		var paste Paste
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

func GetAllPastes() []Paste {
	return getPastes("",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"BurnAfter",
		})
}

func GetPublicPastes() []Paste {
	return getPastes(
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

func GetPublicRedirects() []Paste {
	return getPastes(
		"where Privacy = 'public' and UrlRedirect != '0'",
		[]string{},
		[]string{
			"PasteName",
			"Expire",
			"Privacy",
			"BurnAfter",
			"Syntax",
		},
	)
}

func GetPasteByName(pasteName string) (Paste, error) {
	CleanUpExpiredPastes()
	pastes := getPastes(
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
			"HashedPassword",
			"CreatedAt",
			"UpdatedAt",
		},
	)
	if len(pastes) == 0 {
		return Paste{}, errors.New("Paste not found")
	}
	return pastes[0], nil
}
