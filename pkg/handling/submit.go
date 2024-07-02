package handling

import (
	"errors"
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

	if err := validateSubmit(&submit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isRedirect := shared.IsContentJustUrl(submit.Text)
	pasteName := database.CreatePasteName(isRedirect)

	c.JSON(http.StatusOK, gin.H{
		"text":       submit.Text,
		"expiration": submit.Expiration,
		"burn":       submit.BurnAfter,
		"syntax":     submit.Syntax,
		"privacy":    submit.Privacy,
		"file":       "file",
		"pasteUrl":   pasteName,
	})

	paste := shared.SubmitToPaste(submit, pasteName, isRedirect)
	database.CreatePaste(paste)
}

func parseSubmitForm(c *gin.Context) (shared.Submit, error) {
	var submit shared.Submit
	submit.Text = c.PostForm("text")
	submit.Expiration = c.PostForm("expiration")
	submit.Password = c.PostForm("password")
	submit.Syntax = c.PostForm("syntax")
	submit.Privacy = c.PostForm("privacy")
	burnInt, err := strconv.Atoi(c.PostForm("burn"))
	if err != nil {
		return shared.Submit{}, errors.New("burn must be an integer")
	}
	submit.BurnAfter = burnInt

	form, err := c.MultipartForm()
	if err != nil {
		return shared.Submit{}, errors.New("form error: " + err.Error())
	}

	submit.Files = form.File["file"]
	return submit, nil
}

func validateSubmit(submit *shared.Submit) error {
	if submit.Text == "" && len(submit.Files) == 0 {
		return errors.New("text or file is required")
	}
	encrypt := (submit.Privacy == "private" || submit.Privacy == "secret")
	if submit.Password == "" && encrypt {
		return errors.New("password is required for private or secret pastes")
	}

	if shared.NotAllowedPrivacy(submit.Privacy) {
		return errors.New("invalid privacy")
	}

	if config.Config.DisableEternalPaste && submit.Expiration == "never" {
		submit.Expiration = "1w"
	}

	if config.Config.MaxContentLength > 0 && len(submit.Text) > config.Config.MaxContentLength {
		return errors.New("content is too long")
	}

	if config.Config.MaxFileSize > 0 && len(submit.Files) > 0 {
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

	if config.Config.MaxEncryptionSize > 0 && encrypt && len(submit.Files) > 0 {
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
