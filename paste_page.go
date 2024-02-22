package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePastePage(c *gin.Context) {
	paste, _ := GetPasteByName(c.Param("pasteName"))
	c.HTML(http.StatusOK, "main.html", gin.H{
		"paste": paste,
	})
}
