package main

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Submit struct {
	Text       string                `form:"text" binding:"required"`
	Expiration string                `form:"expiration,omitempty"`
	BurnAfter  int                   `form:"burn,omitempty"`
	Password   string                `form:"password,omitempty"`
	Syntax     string                `form:"syntax,omitempty"`
	Privacy    string                `form:"privacy,omitempty"`
	File       *multipart.FileHeader `form:"file,omitempty"`
}

func HandleSubmit(c *gin.Context) {
	var submit Submit
	if err := c.ShouldBind(&submit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if submit.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}
	pasteName := CreatePasteName()
	hashedPassword := HashPassword(submit.Password)
	c.JSON(http.StatusOK, gin.H{
		"text":       submit.Text,
		"expiration": submit.Expiration,
		"burn":       submit.BurnAfter,
		"password":   hashedPassword,
		"syntax":     submit.Syntax,
		"privacy":    submit.Privacy,
		"file":       submit.File,
		"pasteUrl":   pasteName,
	})

	paste := SubmitToPaste(submit, pasteName, hashedPassword)
	CreatePaste(paste)
}
