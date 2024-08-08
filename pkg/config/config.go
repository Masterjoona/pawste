package config

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Vars ConfigEnv
var Logger *zap.SugaredLogger

const (
	PawsteVersion = ""
	envPrefix     = "PAWSTE_"
)

func (ConfigEnv) InitConfig() {
	Vars = ConfigEnv{
		Salt:                  getEnv("SALT", "banana"),
		Port:                  getEnv("PORT", ":9454"),
		DataDir:               getEnv("DATA_DIR", "pawste_data/"),
		AdminPassword:         getEnv("ADMIN_PASSWORD", "admin"),
		PublicList:            getEnv("PUBLIC_LIST", "true") == "true",
		FileUpload:            getEnv("FILE_UPLOAD", "true") == "true",
		MaxFileSize:           getEnvInt("MAX_FILE_SIZE", "1024 * 1024 * 10"),
		MaxEncryptionSize:     getEnvInt("MAX_ENCRYPTION_SIZE", "1024 * 1024 * 10"),
		MaxContentLength:      getEnvInt("MAX_CONTENT_LENGTH", "5000"),
		FileUploadingPassword: getEnv("FILE_UPLOADING_PASSWORD", ""),
		EternalPaste:          getEnv("ETERNAL_PASTE", "false") == "true",
		MaxExpiryTime:         getEnv("MAX_EXPIRY_TIME", "1w"),
		ReadCount:             getEnv("READ_COUNT", "true") == "true",
		BurnAfter:             getEnv("BURN_AFTER", "true") == "true",
		DefaultExpiry:         getEnv("DEFAULT_EXPIRY", "1w"),
		ShortPasteNames:       getEnv("SHORT_PASTE_NAMES", "false") == "true",
		ShortenRedirectPastes: getEnv("SHORTEN_REDIRECT_PASTES", "false") == "true",
		CountFileUsage:        getEnv("COUNT_FILE_USAGE", "true") == "true",
		AnimeGirlMode:         getEnv("ANIME_GIRL_MODE", "false") == "true",
		LogLevel:              getEnv("LOG_LEVEL", "info"),
	}

	if _, err := os.Stat(Vars.DataDir); os.IsNotExist(err) {
		os.Mkdir(Vars.DataDir, 0755)
	}
	parseLogger()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(envPrefix + key); exists {
		println("Using environment variable", envPrefix+key, "with value", value)
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback string) int {
	if value, exists := os.LookupEnv(envPrefix + key); exists {
		println("Using environment variable", envPrefix+key, "with value", value)
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
	if password == "" || password != Vars.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	Vars.InitConfig()
	c.JSON(http.StatusOK, gin.H{"message": "reloaded config"})
}

func ParseDuration(input string) time.Duration {
	matches := TimeRegex.FindStringSubmatch(input)
	if len(matches) != 3 {
		Logger.Debug("Invalid duration:", input)
		return time.Duration(OneWeek)
	}

	quantity, err := strconv.Atoi(matches[1])
	if err != nil {
		Logger.Debug("Invalid duration:", input)
		return time.Duration(OneWeek)
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
		Logger.Debug("Invalid duration:", input)
		return time.Duration(OneWeek)
	}

	return time.Duration(quantity*multiplier) * time.Second
}

func parseLogger() {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(Vars.LogLevel))
	if err != nil {
		println("Invalid log level:", Vars.LogLevel)
		logLevel = zap.InfoLevel
	}
	Logger = prettyconsole.NewLogger(logLevel).Sugar()
}
