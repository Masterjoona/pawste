package handling

import (
	"net/http"
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/gin-gonic/gin"
)

func HandlePasteRaw(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.String(http.StatusNotFound, "Paste not found")
		return
	}

	reqPassword := c.Request.Header.Get("password")
	needsAuth := paste.NeedsAuth == 1
	if needsAuth && !isValidPassword(reqPassword, paste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	if needsAuth {
		paste.Content = paste.DecryptText(reqPassword)
	}

	database.UpdateReadCount(pasteName)
	c.String(http.StatusOK, paste.Content)
}

func HandleFile(c *gin.Context) {
	queriedPaste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	reqPassword := c.Request.Header.Get("password")
	encrypted := queriedPaste.NeedsAuth == 1 && queriedPaste.Privacy != "readonly"
	if encrypted && !isValidPassword(reqPassword, queriedPaste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	fileDb, err := database.GetFile(queriedPaste.PasteName, c.Param("fileName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	filePath := config.Vars.DataDir + "/" + queriedPaste.PasteName + "/" + fileDb.Name
	if encrypted {
		fileBlob, err := os.ReadFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}
		fileBytes, err := paste.Decrypt(reqPassword, fileBlob)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt file"})
			return
		}
		c.Data(http.StatusOK, fileDb.ContentType, fileBytes)
		return
	}
	if config.Vars.CountFileUsage {
		database.UpdateReadCount(queriedPaste.PasteName)
	}
	c.File(filePath)
}
