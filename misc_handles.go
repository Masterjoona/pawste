package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandlePastePage(c *gin.Context) {
	CleanUpExpiredPastes()
	if ErrorOnInvalidParam(c, "pasteName") {
		return
	}
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	c.HTML(http.StatusOK, "main.html", gin.H{
		"paste": paste,
	})
}

func HandleListPage(c *gin.Context) {
	CleanUpExpiredPastes()
	pastes := GetPublicPastes()
	c.HTML(http.StatusOK, "list.html", gin.H{
		"Pastes": pastes,
	})
}

func HandleGuidePage(c *gin.Context) {
	c.HTML(http.StatusOK, "guide.html", nil)
}

func Redirect(c *gin.Context) {
	CleanUpExpiredPastes()
	if ErrorOnInvalidParam(c, "pasteName") {
		return
	}
	paste, err := GetPasteByName(c.Param("pasteName"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	content := paste.Content
	if strings.HasPrefix(content, "http") || strings.HasPrefix(content, "magnet") {
		c.Redirect(http.StatusFound, content)
		return
	} else {
		c.Redirect(http.StatusFound, "/p/"+c.Param("pasteName"))
		return
	}
}
