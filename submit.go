package main

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	submit.Text = c.PostForm("text")
	submit.Expiration = c.PostForm("expiration")
	submit.Password = c.PostForm("password")
	submit.Syntax = c.PostForm("syntax")
	submit.Privacy = c.PostForm("privacy")
	burnInt, err := strconv.Atoi(c.PostForm("burn"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "burn must be an integer"})
		return
	}
	submit.BurnAfter = burnInt

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form error: " + err.Error()})
		return
	}

	files := form.File["files"]

	if submit.Text == "" && submit.Files == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text or file is required"})
		return
	}
	submit.Files = files

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
