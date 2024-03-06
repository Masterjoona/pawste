package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	UpdateReadCount(c.Param("pasteName"))
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.HTML(http.StatusNotFound, "main.html", gin.H{
			"NotFound": true,
			"Config":   Config,
		})
		return
	}
	c.HTML(http.StatusOK, "main.html", gin.H{
		"Paste":  paste,
		"Config": Config,
	})
}

func HandlePasteJSON(c *gin.Context) {
	UpdateReadCount(c.Param("pasteName"))
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	c.JSON(http.StatusOK, paste)
}

func HandleUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func RedirectHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func HandleRaw(c *gin.Context) {
	UpdateReadCount(c.Param("pasteName"))
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.String(http.StatusNotFound, "Paste not found")
		return
	}
	c.String(http.StatusOK, paste.Content)
}

func Redirect(c *gin.Context) {
	UpdateReadCount(c.Param("pasteName"))
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if paste.UrlRedirect == 0 {
		c.Redirect(http.StatusFound, "/p/"+paste.PasteName)
		return
	}
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

func ReloadConfig(c *gin.Context) {
	var password SubmittedPassword
	if err := c.Bind(&password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if password.Password != Config.AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	InitConfig()
	c.JSON(http.StatusOK, gin.H{"message": "reloaded config"})
}
