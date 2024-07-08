package handling

import (
	"net/http"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
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

	if isEncrypted && password == "" {
		c.Redirect(http.StatusFound, "/p/"+pasteName+"/auth")
		return
	}

	if isEncrypted && paste.Password != database.HashPassword(password) {
		c.Redirect(http.StatusFound, "/p/"+pasteName+"/auth")
		return
	}

	if isEncrypted {
		println("decrypting")
		println(paste.Content)
		paste.Content = paste.DecryptText(password)
		println(paste.Content)
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
	database.UpdateReadCount(pasteName)
	c.JSON(http.StatusOK, paste)
}

func HandleUpdate(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	var newPaste utils.Submit
	if err := c.Bind(&newPaste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if paste.Password != "" {
		if paste.Password != database.HashPassword(newPaste.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}
	}
	paste.Content = newPaste.Text
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
		c.HTML(http.StatusNotFound, "main.html", gin.H{
			"NotFound":      true,
			"config.Config": config.Config,
		})
		return
	}
	c.HTML(http.StatusOK, "main.html", gin.H{
		"Edit":          paste,
		"config.Config": config.Config,
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
	file, err := database.GetFile(paste.PasteName, c.Param("fileName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	c.JSON(http.StatusOK, file)
}

func HandleFile(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	fileDb, err := database.GetFile(paste.PasteName, c.Param("fileName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	if config.Config.CountFileUsage {
		database.UpdateReadCount(paste.PasteName)
	}
	c.File(config.Config.DataDir + "/" + paste.PasteName + "/" + fileDb.Name)
}
