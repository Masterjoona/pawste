package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGuidePage(c *gin.Context) {
	c.HTML(http.StatusOK, "guide.html", nil)
}
