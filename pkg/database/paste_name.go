package database

import (
	"strings"

	"github.com/Masterjoona/pawste/pkg/config"
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

func createPasteName(shorten int) string {
	if config.Vars.ShortPasteNames || (shorten == 1 && config.Vars.ShortenRedirectPastes) {
		return createShortPasteName()
	}
	for {
		var name strings.Builder
		for i := 0; i < 3; i++ {
			name.WriteString(AnimalNames[config.RandomSource.Intn(len(AnimalNames))])
			name.WriteString("-")
		}
		trimmedName := name.String()[:name.Len()-1]
		if !pasteExists(trimmedName) {
			return trimmedName
		}
	}
}

func createShortPasteName() string {
	for {
		var name strings.Builder
		for i := 0; i < 6; i++ {
			name.WriteByte(characters[config.RandomSource.Intn(len(characters))])
		}
		trimmedName := name.String()
		if !pasteExists(trimmedName) {
			return trimmedName
		}
	}
}

func createShortFileName(pasteName string) string {
	for {
		var name strings.Builder
		for i := 0; i < 6; i++ {
			name.WriteByte(characters[config.RandomSource.Intn(len(characters))])
		}
		trimmedName := name.String()
		if !fileExists(pasteName, trimmedName) {
			return trimmedName
		}
	}
}
