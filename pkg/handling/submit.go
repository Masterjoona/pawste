package handling

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
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

	isRedirect := utils.IsContentJustUrl(submit.Text)
	pasteName := database.CreatePasteName(isRedirect)

	paste := utils.SubmitToPaste(submit, pasteName, isRedirect)
	err = database.CreatePaste(paste)
	if err != nil {
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pasteName": pasteName,
	})
}

func parseSubmitForm(c *gin.Context) (utils.Submit, error) {
	var submit utils.Submit
	submit.Text = c.PostForm("content")
	submit.Expiration = c.PostForm("expiration")
	submit.Password = c.PostForm("password")
	submit.Syntax = c.PostForm("syntax")
	submit.Privacy = c.PostForm("privacy")
	burnAfterInt, err := strconv.Atoi(c.PostForm("burnafter"))
	if err != nil {
		return utils.Submit{}, errors.New("burnafter must be an integer")
	}
	submit.BurnAfter = burnAfterInt

	form, err := c.MultipartForm()
	if err != nil {
		return utils.Submit{}, errors.New("form error: " + err.Error())
	}

	submit.Files = form.File["files[]"]
	return submit, nil
}

func validateSubmit(submit *utils.Submit) error {
	hasFiles := len(submit.Files) > 0
	if submit.Text == "" && !hasFiles {
		return errors.New("text or file is required")
	}
	needsAuth := (submit.Privacy == "private" || submit.Privacy == "secret" || submit.Privacy == "readonly")
	if submit.Password == "" && needsAuth {
		return errors.New("password is required for private, secret or readonly pastes")
	}

	if !utils.AllowedOption(submit.Privacy, paste.PrivacyOptions) {
		return errors.New("invalid privacy")
	}

	if !config.Config.EternalPaste && submit.Expiration == "never" {
		submit.Expiration = "1w"
	}

	if config.Config.MaxContentLength > 0 && len(submit.Text) > config.Config.MaxContentLength {
		return errors.New("content is too long")
	}

	maxSizeFiles := utils.TernaryInt((needsAuth && submit.Privacy != "readonly"), config.Config.MaxEncryptionSize, config.Config.MaxFileSize)

	if hasFiles {
		totalSize := 0
		for _, file := range submit.Files {
			if file == nil {
				continue
			}
			totalSize += int(file.Size)
		}
		if totalSize > maxSizeFiles {
			return errors.New("files are too large")
		}
	}

	return nil
}
