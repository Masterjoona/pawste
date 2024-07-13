package handling

import (
	"net/http"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/gin-gonic/gin"
	"github.com/nichady/golte"
)

func RedirectHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
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

func HandleNewPage(c *gin.Context) {
	golte.RenderPage(c.Writer, c.Request, "page/new", map[string]any{
		"fileUpload":        config.Vars.FileUpload,
		"maxFileSize":       config.Vars.MaxFileSize,
		"MaxEncryptionSize": config.Vars.MaxEncryptionSize,
		"maxContentLength":  config.Vars.MaxContentLength,
	})
}

func HandlePaste(c *gin.Context) {
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
