package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func HandleListPage(c *gin.Context) {
	CleanUpExpiredPastes()
	pastes := GetPublicPastes()
	c.HTML(http.StatusOK, "list.html", gin.H{
		"Pastes": pastes,
	})
}
