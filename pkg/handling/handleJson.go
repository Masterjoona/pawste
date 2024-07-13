package handling

import (
	"net/http"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/romana/rlog"
)

func HandlePasteJson(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	reqPassword := c.Request.Header.Get("password")
	needsAuth := paste.NeedsAuth == 1
	if needsAuth && !isValidPassword(reqPassword, paste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	if needsAuth && paste.Privacy != "readonly" {
		paste.Content = paste.DecryptText(reqPassword)
	}

	paste.Files = database.GetFiles(pasteName)

	database.UpdateReadCount(pasteName)
	c.JSON(http.StatusOK, paste)
}

func HandleEditJson(c *gin.Context) {
	queriedPaste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	pasteFiles := database.GetFiles(queriedPaste.PasteName)

	var newPaste utils.PasteUpdate
	if err := c.Bind(&newPaste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form error: " + err.Error()})
		return
	}

	newPaste.FilesMultiPart = form.File["files[]"]

	needsAuth := queriedPaste.Privacy == "private" || queriedPaste.Privacy == "secret"
	if needsAuth && queriedPaste.Password != database.HashPassword(newPaste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	if config.Vars.MaxContentLength < len(newPaste.Content) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content too long"})
		return
	}

	currentFileSizeTotal := 0
	for _, file := range pasteFiles {
		currentFileSizeTotal += file.Size
	}
	filesToBeRemoved := newPaste.RemovedFiles
	for _, file := range pasteFiles {
		if filesToBeRemoved == nil {
			break
		}
		for _, fileName := range filesToBeRemoved {
			if fileName == file.Name {
				currentFileSizeTotal -= file.Size
			}
		}
	}
	for _, file := range newPaste.FilesMultiPart {
		currentFileSizeTotal += int(file.Size)
		fileName, fileSize, fileBlob := utils.ConvertMultipartFile(file)
		newPaste.Files = append(newPaste.Files, paste.File{
			Name:        fileName,
			Size:        fileSize,
			Blob:        fileBlob,
			ContentType: file.Header.Get("Content-Type"),
		})
	}

	if config.Vars.MaxFileSize > 0 && currentFileSizeTotal > config.Vars.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size too big"})
		return
	}

	err = database.UpdatePaste(queriedPaste.PasteName, newPaste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		rlog.Errorf("Failed to update paste: %s", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "paste updated"})
}

func HandlePasteDelete(c *gin.Context) {
	pasteName := c.Param("pasteName")
	queriedPaste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	if queriedPaste.NeedsAuth == 1 {
		password := c.Request.Header.Get("password")
		if password == "" || !isValidPassword(password, queriedPaste.Password) && password != config.Vars.AdminPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}
	}

	if err := database.DeletePaste(pasteName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "paste deleted"})
}

func HandleFileJson(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	reqPassword := c.Request.Header.Get("password")
	needsAuth := paste.NeedsAuth == 1
	if needsAuth && paste.Privacy != "readonly" && !isValidPassword(reqPassword, paste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	file, err := database.GetFile(paste.PasteName, c.Param("fileName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	c.JSON(http.StatusOK, file)
}

func HandleAdminJson(c *gin.Context) {
	passwd := c.Request.Header.Get("password")
	if passwd == "" || passwd != config.Vars.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"config": config.Vars, "pastes": database.GetAllPastes()})
}
