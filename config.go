package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/romana/rlog"
)

type ConfigEnv struct {
	Salt          string
	Port          string
	DataDir       string
	AdminPassword string

	PublicList bool
	PublicURL  string

	DefaultExpiryTime   string
	NoFileUpload        bool
	MaxFileSize         int
	MaxEncryptionSize   int
	MaxContentLength    int
	UploadingPassword   string
	DisableEternalPaste bool
	DisableReadCount    bool
	DisableBurnAfter    bool
	DefaultExpiry       string
	ShortPasteNames     bool
	IUnderstandTheRisks bool
}

var Config ConfigEnv

const (
	PawsteVersion = ""
	envPrefix     = "PAWSTE_"
)

func InitConfig() {
	Config = ConfigEnv{
		Salt:                getEnv("SALT", "banana"),
		Port:                getEnv("PORT", ":9454"),
		DataDir:             getEnv("DATA_DIR", "pawste_data/"),
		AdminPassword:       getEnv("ADMIN_PASSWORD", "admin"),
		PublicList:          getEnv("PUBLIC_LIST", "true") == "true",
		PublicURL:           getEnv("PUBLIC_URL", "http://localhost:"+getEnv("PORT", ":9454")),
		NoFileUpload:        getEnv("NO_FILE_UPLOAD", "false") == "true",
		MaxFileSize:         getEnvInt("MAX_FILE_SIZE", "1024 * 1024 * 10"),
		MaxEncryptionSize:   getEnvInt("MAX_ENCRYPTION_SIZE", "1024 * 1024 * 10"),
		MaxContentLength:    getEnvInt("MAX_CONTENT_LENGTH", "1024 * 1024"),
		UploadingPassword:   getEnv("UPLOADING_PASSWORD", ""),
		DisableEternalPaste: getEnv("DISABLE_ETERNAL_PASTE", "false") == "true",
		DisableReadCount:    getEnv("DISABLE_READ_COUNT", "false") == "true",
		DisableBurnAfter:    getEnv("DISABLE_BURN_AFTER", "false") == "true",
		DefaultExpiry:       getEnv("DEFAULT_EXPIRY", "1w"),
		ShortPasteNames:     getEnv("SHORT_PASTE_NAMES", "false") == "true",
		IUnderstandTheRisks: getEnv("I_UNDERSTAND_THE_RISKS", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(envPrefix + key); exists {
		rlog.Info("Using environment variable", envPrefix+key, "with value", value)
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback string) int {
	if value, exists := os.LookupEnv(envPrefix + key); exists {
		return CalculateIntFromString(value)
	}
	return CalculateIntFromString(fallback)
}

func CalculateIntFromString(s string) int {
	parts := strings.Split(s, "*")
	result := 1
	for _, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			panic(err)
		}
		result *= num
	}
	return result
}
