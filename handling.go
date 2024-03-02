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
		c.Redirect(http.StatusNotFound, "/")
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

type PasteLists struct {
	Pastes    []Paste
	Redirects []Paste
}

func ListHandler() interface{} {
	return PasteLists{
		Pastes:    GetPublicPastes(),
		Redirects: GetPublicRedirects(),
	}

}
