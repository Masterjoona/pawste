package handling

import (
	"net/http"

	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/shared"
	"github.com/Masterjoona/pawste/pkg/shared/config"
	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
	"github.com/romana/rlog"
)

func HandlePage(
	settings map[string]interface{},
	function func() interface{},
	value string,
) gin.HandlerFunc {
	settings["config.Config"] = config.Config
	return func(c *gin.Context) {
		if function != nil {
			settings[value] = function()
		}
		c.HTML(http.StatusOK, "main.html", settings)
	}
}

func HandlePastePage(c *gin.Context) {
	pasteName := c.Param("pasteName")
	paste, err := database.GetPasteByName(pasteName)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	golte.RenderPage(c.Writer, c.Request, "page/paste", map[string]any{
		"paste": paste,
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
	database.UpdateReadCount(pasteName)
	c.JSON(http.StatusOK, paste)
}

func HandleUpdate(c *gin.Context) {
	paste, err := database.GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	var newPaste shared.Submit
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
	paste.UpdatedAt = shared.GetCurrentDate()
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

func HandleRaw(c *gin.Context) {
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

func HandleFile(c *gin.Context) {
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

// funny
func AdminHandler() interface{} {
	return database.GetAllPastes()
}

func ListHandler() interface{} {
	return paste.PasteLists{
		Pastes:    database.GetPublicPastes(),
		Redirects: database.GetPublicRedirects(),
	}
}
