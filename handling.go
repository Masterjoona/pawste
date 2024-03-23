package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romana/rlog"
)

func HandlePage(
	settings map[string]interface{},
	function func() interface{},
	value string,
) gin.HandlerFunc {
	CleanUpExpiredPastes()
	settings["Config"] = Config
	return func(c *gin.Context) {
		if function != nil {
			settings[value] = function()
		}
		c.HTML(http.StatusOK, "main.html", settings)
	}
}

func HandlePastePage(c *gin.Context) {
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.HTML(http.StatusNotFound, "main.html", gin.H{
			"NotFound": true,
			"Config":   Config,
		})
		return
	}
	UpdateReadCount(c.Param("pasteName"))
	c.HTML(http.StatusOK, "main.html", gin.H{
		"Paste":  paste,
		"Config": Config,
	})
}

func HandlePasteJSON(c *gin.Context) {
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	UpdateReadCount(c.Param("pasteName"))
	c.JSON(http.StatusOK, paste)
}

func HandleUpdate(c *gin.Context) {
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	var newPaste Submit
	if err := c.Bind(&newPaste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if paste.HashedPassword != "" {
		if paste.HashedPassword != HashPassword(newPaste.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}
	}
	paste.Content = newPaste.Text
	paste.UpdatedAt = GetCurrentDate()
	err = UpdatePaste(paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		rlog.Errorf("Failed to update paste: %s", err)
		return
	}
	c.JSON(http.StatusOK, paste)
}

func HandleEdit(c *gin.Context) {
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.HTML(http.StatusNotFound, "main.html", gin.H{
			"NotFound": true,
			"Config":   Config,
		})
		return
	}
	c.HTML(http.StatusOK, "main.html", gin.H{
		"Edit":   paste,
		"Config": Config,
	})
}

func RedirectHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func HandleRaw(c *gin.Context) {
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.String(http.StatusNotFound, "Paste not found")
		return
	}
	UpdateReadCount(c.Param("pasteName"))
	c.String(http.StatusOK, paste.Content)
}

func Redirect(c *gin.Context) {
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if paste.UrlRedirect == 0 {
		c.Redirect(http.StatusFound, "/p/"+paste.PasteName)
		return
	}
	UpdateReadCount(c.Param("pasteName"))
	c.Redirect(http.StatusFound, paste.Content)

}

// funny
func AdminHandler() interface{} {
	return GetAllPastes()
}

func ListHandler() interface{} {
	return PasteLists{
		Pastes:    GetPublicPastes(),
		Redirects: GetPublicRedirects(),
	}
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
