package config

import (
	"regexp"
	"time"
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
	UploadingPassword     string
	EternalPaste          bool
	MaxExpiryTime         int
	ReadCount             bool
	BurnAfter             bool
	DefaultExpiry         string
	ShortPasteNames       bool
	ShortenRedirectPastes bool
	CountFileUsage        bool
	AnimeGirlMode         bool
}

var TimeRegex = regexp.MustCompile(`^(\d+)([smhdwMy])$`)

const OneWeek = time.Hour * 24 * 7
