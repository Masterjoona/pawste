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
	CleanUpExpiredPastes()
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.HTML(http.StatusNotFound, "main.html", gin.H{
			"404":    true,
			"Config": Config,
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
	CleanUpExpiredPastes()
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.String(http.StatusNotFound, "Paste not found")
		return
	}
	c.String(http.StatusOK, paste.Content)
}

func Redirect(c *gin.Context) {
	CleanUpExpiredPastes()
	//paste, err := GetPasteByName(c.Param("pasteName"))
	// refactor
}

// funny
func adminHandler() interface{} {
	return GetAllPastes()
}

func listHandler() interface{} {
	return GetPublicPastes()
}
