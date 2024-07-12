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
	queriedPaste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	needsAuth := queriedPaste.NeedsAuth == 1
	burnAfter := queriedPaste.BurnAfter == 1
	if needsAuth && queriedPaste.Privacy != "readonly" {
		golte.RenderPage(c.Writer, c.Request, "page/paste", map[string]any{
			"needsAuth": true,
			"paste":     paste.Paste{},
			"burnAfter": burnAfter,
		})
		return
	}

	if burnAfter && c.Query("read") == "" {
		golte.RenderPage(c.Writer, c.Request, "page/oneview", nil)
		return
	}
	queriedPaste.Files = database.GetFiles(pasteName)
	golte.RenderPage(c.Writer, c.Request, "page/paste", map[string]any{
		"paste":     queriedPaste,
		"needsAuth": needsAuth,
	})
	database.UpdateReadCount(pasteName)
}

func HandlePasteJSON(c *gin.Context) {
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

func HandleUpdate(c *gin.Context) {
	queriedPaste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	pasteFiles := database.GetFiles(queriedPaste.PasteName)
	needsAuth := queriedPaste.Privacy == "private" || queriedPaste.Privacy == "secret"
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

	if needsAuth && queriedPaste.Password != database.HashPassword(newPaste.Password) {
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

	if config.Config.MaxFileSize > 0 && currentFileSizeTotal > config.Config.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size too big"})
		return
	}

	queriedPaste.Content = newPaste.Content
	err = database.UpdatePaste(queriedPaste.PasteName, newPaste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		rlog.Errorf("Failed to update paste: %s", err)
		return
	}
	c.JSON(http.StatusOK, queriedPaste)
}

func HandleEdit(c *gin.Context) {
	queriedPaste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if queriedPaste.NeedsAuth == 1 {
		tmpPaste := paste.Paste{}
		tmpPaste.Files = []paste.File{}
		golte.RenderPage(c.Writer, c.Request, "page/edit", map[string]any{
			"needsAuth": queriedPaste.NeedsAuth == 1,
			"paste":     tmpPaste,
		})
		return
	}
	queriedPaste.Files = database.GetFiles(queriedPaste.PasteName)
	golte.RenderPage(c.Writer, c.Request, "page/edit", map[string]any{
		"paste": queriedPaste,
	})

}

func RedirectHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func HandleList(c *gin.Context) {
	golte.RenderPage(c.Writer, c.Request, "page/list", map[string]any{
		"pastes": database.GetAllPublicPastes(),
	})
}

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

func Redirect(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if paste.UrlRedirect == 0 {
		c.Redirect(http.StatusFound, "/p/"+pasteName)
		return
	}
	database.UpdateReadCount(pasteName)
	c.Redirect(http.StatusFound, paste.Content)
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
	filePath := config.Config.DataDir + "/" + queriedPaste.PasteName + "/" + fileDb.Name
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
	if config.Config.CountFileUsage {
		database.UpdateReadCount(queriedPaste.PasteName)
	}
	c.File(filePath)
}
func HandlePasteDelete(c *gin.Context) {
	pasteName := c.Param("pasteName")
	queriedPaste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	if queriedPaste.NeedsAuth == 1 {
		var deletePasswd config.PasswordJSON
		if err := c.ShouldBindJSON(&deletePasswd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
			return
		}

		if !isValidPassword(deletePasswd.Password, queriedPaste.Password) && deletePasswd.Password != config.Config.AdminPassword {
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

func HandleAdmin(c *gin.Context) {
	golte.RenderPage(c.Writer, c.Request, "page/admin", map[string]any{
		"config": config.ConfigEnv{},
		"pastes": []paste.Paste{},
	})
}

func HandleAdminJSON(c *gin.Context) {
	passwd := c.Request.Header.Get("password")
	if passwd == "" || passwd != config.Config.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"config": config.Config, "pastes": database.GetAllPastes()})
}
