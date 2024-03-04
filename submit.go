package main

import (
	"errors"
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
	submit, err := parseSubmitForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateSubmit(submit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pasteName := CreatePasteName()

	c.JSON(http.StatusOK, gin.H{
		"text":       submit.Text,
		"expiration": submit.Expiration,
		"burn":       submit.BurnAfter,
		"syntax":     submit.Syntax,
		"privacy":    submit.Privacy,
		"file":       "file",
		"pasteUrl":   pasteName,
	})

	paste := SubmitToPaste(submit, pasteName)
	CreatePaste(paste)
}

func parseSubmitForm(c *gin.Context) (Submit, error) {
	var submit Submit
	submit.Text = c.PostForm("text")
	submit.Expiration = c.PostForm("expiration")
	submit.Password = c.PostForm("password")
	submit.Syntax = c.PostForm("syntax")
	submit.Privacy = c.PostForm("privacy")
	burnInt, err := strconv.Atoi(c.PostForm("burn"))
	if err != nil {
		return Submit{}, errors.New("burn must be an integer")
	}
	submit.BurnAfter = burnInt

	form, err := c.MultipartForm()
	if err != nil {
		return Submit{}, errors.New("form error: " + err.Error())
	}

	submit.Files = form.File["files"]
	return submit, nil
}

func validateSubmit(submit Submit) error {
	if submit.Text == "" && len(submit.Files) == 0 {
		return errors.New("text or file is required")
	}

	if submit.Password == "" && (submit.Privacy == "private" || submit.Privacy == "secret") {
		return errors.New("password is required for private or secret pastes")
	}
	return nil
}
