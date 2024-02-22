package main

import (
	"math/rand"
	"time"
)

var AnimalNames = []string{"ant", "eel", "mole", "sloth", "ape", "emu", "monkey", "snail", "bat", "falcon", "mouse",
	"snake", "bear", "fish", "otter", "spider", "bee", "fly", "parrot", "squid", "bird", "fox",
	"panda", "swan", "bison", "frog", "pig", "tiger", "camel", "gecko", "pigeon", "toad", "cat",
	"goat", "pony", "turkey", "cobra", "goose", "pug", "turtle", "crow", "hawk", "rabbit", "viper",
	"deer", "horse", "rat", "wasp", "dog", "jaguar", "raven", "whale", "dove", "koala", "seal",
	"wolf", "duck", "lion", "shark", "worm", "eagle", "lizard", "sheep", "zebra",
}

func CreatePasteName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var concatted string
	for i := 0; i < 3; i++ {
		concatted += AnimalNames[r.Intn(len(AnimalNames))] + "-"
	}
	return concatted[:len(concatted)-1]
}
