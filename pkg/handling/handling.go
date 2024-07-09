package handling

import (
	"net/http"
	"os"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
	"github.com/romana/rlog"
)

func HandlePastePage(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	password := c.PostForm("password")
	isEncrypted := paste.Privacy == "private" || paste.Privacy == "secret"

	if isEncrypted && (password == "" || paste.Password != database.HashPassword(password)) {
		c.Redirect(http.StatusFound, "/p/"+pasteName+"/auth")
		return
	}

	if isEncrypted {
		paste.Content = paste.DecryptText(password)
	}

	if paste.BurnAfter == 1 && c.Query("read") == "" {
		golte.RenderPage(c.Writer, c.Request, "page/oneview", map[string]any{
			"isEncrypted": isEncrypted,
			"password":    password,
		})
		return
	}

	golte.RenderPage(c.Writer, c.Request, "page/paste", map[string]any{
		"paste": paste,
		"files": database.GetFiles(pasteName),
	})
	database.UpdateReadCount(pasteName)
}

func HandlePastePostAuth(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if paste.Password != database.HashPassword(c.PostForm("password")) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"password": "correct"})
}

func HandlePasteJSON(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	reqPassword := c.Request.Header.Get("password")
	isEncrypted := (paste.Privacy == "private" || paste.Privacy == "secret")
	if verifyAccess(isEncrypted, reqPassword, paste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	if isEncrypted {
		paste.Content = paste.DecryptText(reqPassword)
	}

	database.UpdateReadCount(pasteName)
	c.JSON(http.StatusOK, paste)
}

func HandleUpdate(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	pasteFiles := database.GetFiles(paste.PasteName)
	isEncrypted := paste.Privacy == "private" || paste.Privacy == "secret"
	var newPaste utils.PasteUpdate
	if err := c.Bind(&newPaste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isEncrypted && paste.Password != database.HashPassword(newPaste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	if config.Config.MaxContentLength > 0 && len(newPaste.Content) > config.Config.MaxContentLength {
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

	for _, file := range newPaste.Files {
		currentFileSizeTotal += int(file.Size)
	}

	if config.Config.MaxFileSize > 0 && currentFileSizeTotal > config.Config.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size too big"})
		return
	}

	paste.Content = newPaste.Content
	paste.UpdatedAt = utils.GetCurrentDate()
	err = database.UpdatePaste(paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		rlog.Errorf("Failed to update paste: %s", err)
		return
	}
	c.JSON(http.StatusOK, paste)
}

func HandleEdit(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	golte.RenderPage(c.Writer, c.Request, "page/edit", map[string]any{
		"paste": paste,
		"files": database.GetFiles(paste.PasteName),
	})

}

func RedirectHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func HandlePasteRaw(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.String(http.StatusNotFound, "Paste not found")
		return
	}

	reqPassword := c.Request.Header.Get("password")
	isEncrypted := (paste.Privacy == "private" || paste.Privacy == "secret")
	if verifyAccess(isEncrypted, reqPassword, paste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	if isEncrypted {
		paste.Content = paste.DecryptText(reqPassword)
	}

	database.UpdateReadCount(pasteName)
	c.String(http.StatusOK, paste.Content)
}

func Redirect(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	database.UpdateReadCount(pasteName)
	if paste.UrlRedirect == 0 {
		c.Redirect(http.StatusFound, "/p/"+pasteName)
		return
	}
	c.Redirect(http.StatusFound, paste.Content)
}

func HandleFileJson(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	reqPassword := c.Request.Header.Get("password")
	isEncrypted := (paste.Privacy == "private" || paste.Privacy == "secret")
	if verifyAccess(isEncrypted, reqPassword, paste.Password) {
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

func HandleFile(c *gin.Context) {
	queriedPaste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	reqPassword := c.Request.Header.Get("password")
	isEncrypted := (queriedPaste.Privacy == "private" || queriedPaste.Privacy == "secret")
	if verifyAccess(isEncrypted, reqPassword, queriedPaste.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	fileDb, err := database.GetFile(queriedPaste.PasteName, c.Param("fileName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	filePath := config.Config.DataDir + "/" + queriedPaste.PasteName + "/" + fileDb.Name
	if isEncrypted {
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
	if config.Config.CountFileUsage {
		database.UpdateReadCount(queriedPaste.PasteName)
	}
	c.File(filePath)

}
