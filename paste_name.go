package main

import (
	"math/rand"
	"strings"
	"time"
)

var AnimalNames = []string{
	"ant",
	"eel",
	"mole",
	"sloth",
	"ape",
	"emu",
	"monkey",
	"snail",
	"bat",
	"falcon",
	"mouse",
	"snake",
	"bear",
	"fish",
	"otter",
	"spider",
	"bee",
	"fly",
	"parrot",
	"squid",
	"bird",
	"fox",
	"panda",
	"swan",
	"bison",
	"frog",
	"pig",
	"tiger",
	"camel",
	"gecko",
	"pigeon",
	"toad",
	"cat",
	"goat",
	"pony",
	"turkey",
	"cobra",
	"goose",
	"pug",
	"turtle",
	"crow",
	"hawk",
	"rabbit",
	"viper",
	"deer",
	"horse",
	"rat",
	"wasp",
	"dog",
	"jaguar",
	"raven",
	"whale",
	"dove",
	"koala",
	"seal",
	"wolf",
	"duck",
	"lion",
	"shark",
	"worm",
	"eagle",
	"lizard",
	"sheep",
	"zebra",
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func CreatePasteName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if Config.ShortPasteNames {
		return createShortPasteName(r)
	}
	for {
		var name strings.Builder
		for i := 0; i < 3; i++ {
			name.WriteString(AnimalNames[r.Intn(len(AnimalNames))])
			name.WriteString("-")
		}
		trimmedName := name.String()[:name.Len()-1]
		if !pasteExists(trimmedName) {
			return trimmedName
		}
	}
}

func createShortPasteName(r *rand.Rand) string {
	for {
		var name strings.Builder
		for i := 0; i < 6; i++ {
			name.WriteByte(characters[r.Intn(len(characters))])
		}
		trimmedName := name.String()
		if !pasteExists(trimmedName) {
			return trimmedName
		}
	}
}

func pasteExists(name string) bool {
	var exists bool
	err := PasteDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pastes WHERE PasteName = ?)", name).
		Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists
}
