package config

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/romana/rlog"
)

var Config ConfigEnv

const (
	PawsteVersion = ""
	envPrefix     = "PAWSTE_"
)

func (ConfigEnv) InitConfig() {
	Config = ConfigEnv{
		Salt:                  getEnv("SALT", "banana"),
		Port:                  getEnv("PORT", ":9454"),
		DataDir:               getEnv("DATA_DIR", "pawste_data/"),
		AdminPassword:         getEnv("ADMIN_PASSWORD", "admin"),
		PublicList:            getEnv("PUBLIC_LIST", "true") == "true",
		FileUpload:            getEnv("FILE_UPLOAD", "true") == "true",
		MaxFileSize:           getEnvInt("MAX_FILE_SIZE", "1024 * 1024 * 10"),
		MaxEncryptionSize:     getEnvInt("MAX_ENCRYPTION_SIZE", "1024 * 1024 * 10"),
		MaxContentLength:      getEnvInt("MAX_CONTENT_LENGTH", "1024 * 1024"),
		UploadingPassword:     getEnv("UPLOADING_PASSWORD", ""),
		EternalPaste:          getEnv("ETERNAL_PASTE", "false") == "true",
		MaxExpiryTime:         ParseDuration(getEnv("MAX_EXPIRY_TIME", "1w")),
		ReadCount:             getEnv("READ_COUNT", "true") == "true",
		BurnAfter:             getEnv("BURN_AFTER", "true") == "true",
		DefaultExpiry:         getEnv("DEFAULT_EXPIRY", "1w"),
		ShortPasteNames:       getEnv("SHORT_PASTE_NAMES", "false") == "true",
		ShortenRedirectPastes: getEnv("SHORTEN_REDIRECT_PASTES", "false") == "true",
		CountFileUsage:        getEnv("COUNT_FILE_USAGE", "true") == "true",
		AnimeGirlMode:         getEnv("ANIME_GIRL_MODE", "false") == "true",
	}

	if _, err := os.Stat(Config.DataDir); os.IsNotExist(err) {
		os.Mkdir(Config.DataDir, 0755)
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
		return calculateIntFromString(value)
	}
	return calculateIntFromString(fallback)
}

func calculateIntFromString(s string) int {
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

func (ConfigEnv) ReloadConfig(c *gin.Context) {
	password := c.Request.Header.Get("password")
	if password == "" || password != Config.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	Config.InitConfig()
	c.JSON(http.StatusOK, gin.H{"message": "reloaded config"})
}

func ParseDuration(input string) int {
	matches := TimeRegex.FindStringSubmatch(input)
	if len(matches) != 3 {
		return int(OneWeek)
	}

	quantity, err := strconv.Atoi(matches[1])
	if err != nil {
		return int(OneWeek)
	}

	unit := matches[2]
	unitMultipliers := map[string]int{
		"s": 1,
		"m": 60,
		"h": 3600,
		"d": 86400,
		"w": 604800,
		"M": 2592000,
	}

	multiplier, exists := unitMultipliers[unit]
	if !exists {
		return int(OneWeek)
	}

	return quantity * multiplier
}
