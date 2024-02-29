package main

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileMultiPart struct {
	File *multipart.FileHeader
}

type Submit struct {
	Text       string                  `form:"text,omitempty"`
	Expiration string                  `form:"expiration,omitempty"`
	BurnAfter  int                     `form:"burn,omitempty"`
	Password   string                  `form:"password,omitempty"`
	Syntax     string                  `form:"syntax,omitempty"`
	Privacy    string                  `form:"privacy,omitempty"`
	Files      []*multipart.FileHeader `form:"file,omitempty"`
}

func HandleSubmit(c *gin.Context) {
	var submit Submit
	if err := c.ShouldBind(&submit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if submit.Text == "" || submit.Files[0] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text or file is required"})
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
		"file":       "file",
		"pasteUrl":   pasteName,
	})

	paste := SubmitToPaste(submit, pasteName, hashedPassword)
	CreatePaste(paste, submit.Password)
	// i feel like im doing something wrong here
}
