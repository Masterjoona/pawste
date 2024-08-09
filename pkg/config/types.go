package config

import (
	"regexp"
	"time"

	"math/rand"
)

type ConfigEnv struct {
	Salt                  string
	Port                  string
	DataDir               string
	AdminPassword         string
	PublicList            bool
	FileUpload            bool
	MaxFileSize           int
	MaxEncryptionSize     int
	MaxContentLength      int
	FileUploadingPassword string
	EternalPaste          bool
	MaxExpiryTime         string
	ReadCount             bool
	BurnAfter             bool
	DefaultExpiry         string
	ShortPasteNames       bool
	ShortenRedirectPastes bool
	CountFileUsage        bool
	AnimeGirlMode         bool
	LogLevel              string
	AnonymiseFileNames    bool
}

var TimeRegex = regexp.MustCompile(`^(\d+)([smhdwMy])$`)

const OneWeek = time.Hour * 24 * 7

var RandomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
