package handling

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/shared"
	"github.com/Masterjoona/pawste/pkg/shared/config"
	"github.com/gin-gonic/gin"
)

func HandleSubmit(c *gin.Context) {
	submit, err := parseSubmitForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = validateSubmit(&submit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isRedirect := shared.IsContentJustUrl(submit.Text)
	pasteName := database.CreatePasteName(isRedirect)

	paste := shared.SubmitToPaste(submit, pasteName, isRedirect)
	err = database.CreatePaste(paste)
	if err != nil {
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pasteName": pasteName,
	})
}

func parseSubmitForm(c *gin.Context) (shared.Submit, error) {
	var submit shared.Submit
	submit.Text = c.PostForm("content")
	submit.Expiration = c.PostForm("expiration")
	submit.Password = c.PostForm("password")
	submit.Syntax = c.PostForm("syntax")
	submit.Privacy = c.PostForm("privacy")
	burnAfterInt, err := strconv.Atoi(c.PostForm("burnafter"))
	if err != nil {
		return shared.Submit{}, errors.New("burnafter must be an integer")
	}
	submit.BurnAfter = burnAfterInt

	form, err := c.MultipartForm()
	if err != nil {
		return shared.Submit{}, errors.New("form error: " + err.Error())
	}

	submit.Files = form.File["files[]"]
	fmt.Println(form.File["files[]"])
	return submit, nil
}

func validateSubmit(submit *shared.Submit) error {
	hasFiles := len(submit.Files) > 0
	if submit.Text == "" && !hasFiles {
		return errors.New("text or file is required")
	}
	encrypt := (submit.Privacy == "private" || submit.Privacy == "secret")
	if submit.Password == "" && encrypt {
		return errors.New("password is required for private or secret pastes")
	}

	if shared.NotAllowedPrivacy(submit.Privacy) {
		return errors.New("invalid privacy")
	}

	if !config.Config.EternalPaste && submit.Expiration == "never" {
		submit.Expiration = "1w"
	}

	if config.Config.MaxContentLength > 0 && len(submit.Text) > config.Config.MaxContentLength {
		return errors.New("content is too long")
	}

	if !encrypt && hasFiles && config.Config.MaxFileSize > 0 {
		totalSize := 0
		for _, file := range submit.Files {
			if file == nil {
				continue
			}
			totalSize += int(file.Size)
		}
		if totalSize > config.Config.MaxFileSize {
			return errors.New("file size is too large")
		}
	}

	if encrypt && hasFiles && config.Config.MaxEncryptionSize > 0 {
		totalSize := 0
		for _, file := range submit.Files {
			if file == nil {
				continue
			}
			totalSize += int(file.Size)
		}
		if totalSize > config.Config.MaxEncryptionSize {
			return errors.New("file size is too large for encryption")
		}
	}

	return nil
}
