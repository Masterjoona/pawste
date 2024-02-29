package main

import (
	"os"
	"strconv"
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
	MaxFileSize         int64
	MaxEncryptionSize   int64
	MaxContentLength    int64
	UploadingPassword   string
	DisableEternalPaste bool
	DisableReadCount    bool
	DisableBurnAfter    bool
	DefaultExpiry       string
	ShortPasteName      bool
}

var (
	Config        ConfigEnv
	PawsteVersion string
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
		MaxFileSize:         int64(getEnvInt("MAX_FILE_SIZE", 1024*1024*10)),
		MaxEncryptionSize:   int64(getEnvInt("MAX_ENCRYPTION_SIZE", 1024*1024*10)),
		MaxContentLength:    int64(getEnvInt("MAX_CONTENT_LENGTH", 1024*1024*10)),
		UploadingPassword:   getEnv("UPLOADING_PASSWORD", ""),
		DisableEternalPaste: getEnv("DISABLE_ETERNAL_PASTE", "false") == "true",
		DisableReadCount:    getEnv("DISABLE_READ_COUNT", "false") == "true",
		DisableBurnAfter:    getEnv("DISABLE_BURN_AFTER", "false") == "true",
		DefaultExpiry:       getEnv("DEFAULT_EXPIRY", "1w"),
		ShortPasteName:      getEnv("SHORT_PASTE_NAME", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(envPrefix + key); exists {
		println("Using env var", envPrefix+key, "with value", value)
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(envPrefix + key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
