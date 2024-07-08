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
	port := getEnv("PORT", ":9454")
	Config = ConfigEnv{
		Salt:                  getEnv("SALT", "banana"),
		Port:                  port,
		DataDir:               getEnv("DATA_DIR", "pawste_data/"),
		AdminPassword:         getEnv("ADMIN_PASSWORD", "admin"),
		PublicList:            getEnv("PUBLIC_LIST", "true") == "true",
		PublicURL:             getEnv("PUBLIC_URL", "http://localhost"+port),
		FileUpload:            getEnv("FILE_UPLOAD", "true") == "true",
		MaxFileSize:           getEnvInt("MAX_FILE_SIZE", "1024 * 1024 * 10"),
		MaxEncryptionSize:     getEnvInt("MAX_ENCRYPTION_SIZE", "1024 * 1024 * 10"),
		MaxContentLength:      getEnvInt("MAX_CONTENT_LENGTH", "1024 * 1024"),
		UploadingPassword:     getEnv("UPLOADING_PASSWORD", ""),
		EternalPaste:          getEnv("ETERNAL_PASTE", "false") == "true",
		ReadCount:             getEnv("READ_COUNT", "true") == "true",
		BurnAfter:             getEnv("BURN_AFTER", "true") == "true",
		DefaultExpiry:         getEnv("DEFAULT_EXPIRY", "1w"),
		ShortPasteNames:       getEnv("SHORT_PASTE_NAMES", "false") == "true",
		ShortenRedirectPastes: getEnv("SHORTEN_REDIRECT_PASTES", "false") == "true",
		CountFileUsage:        getEnv("COUNT_FILE_USAGE", "true") == "true",
		IUnderstandTheRisks:   getEnv("I_UNDERSTAND_THE_RISKS", "false") == "true",
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

type PasswordSubmission struct {
	Password string `json:"password"`
}

func (ConfigEnv) ReloadConfig(c *gin.Context) {
	var password PasswordSubmission
	if err := c.Bind(&password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "give a password dumbass"})
		return
	}
	if password.Password != Config.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	Config.InitConfig()
	c.JSON(http.StatusOK, gin.H{"message": "reloaded config"})
}
