package handling

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/database"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/romana/rlog"
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

	paste, err := submitToPaste(submit, pasteName, isRedirect)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	err = database.CreatePaste(paste)
	if err != nil {
		rlog.Error("Error creating paste", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	for i := range paste.Files {
		paste.Files[i].Blob = nil
	}

	c.JSON(http.StatusOK,
		paste,
	)
}

func submitToPaste(submit Submit, pasteName string, isRedirect int) (paste.Paste, error) {
	var files []paste.File
	for _, file := range submit.Files {
		if file == nil {
			continue
		}
		fileName, fileSize, fileBlob, err := convertMultipartFile(file)
		if err != nil {
			return paste.Paste{}, err
		}
		files = append(files, paste.File{
			Name:        fileName,
			Size:        fileSize,
			Blob:        fileBlob,
			ContentType: file.Header.Get("Content-Type"),
		})
	}
	todaysDate := time.Now().Unix()
	return paste.Paste{
		PasteName:   pasteName,
		Expire:      utils.HumanTimeToUnix(submit.Expiration),
		Privacy:     submit.Privacy,
		NeedsAuth:   utils.Ternary((submit.Password == ""), 0, 1),
		ReadCount:   0,
		ReadLast:    todaysDate,
		BurnAfter:   utils.Ternary(config.Vars.BurnAfter, submit.BurnAfter, 0),
		Content:     submit.Text,
		Syntax:      submit.Syntax,
		Password:    submit.Password,
		Files:       files,
		UrlRedirect: isRedirect,
		CreatedAt:   todaysDate,
		UpdatedAt:   todaysDate,
	}, nil
}

func convertMultipartFile(file *multipart.FileHeader) (string, int, []byte, error) {
	src, err := file.Open()
	if err != nil {
		rlog.Error("Could not open multipart file", err)
		return "", 0, nil, err
	}
	defer src.Close()

	fileBlob, err := io.ReadAll(src)
	if err != nil {
		rlog.Error("Could not read multipart file", err)
		return "", 0, nil, err
	}
	return file.Filename, len(fileBlob), fileBlob, nil
}

func parseSubmitForm(c *gin.Context) (Submit, error) {
	var submit Submit
	submit.Text = c.PostForm("content")
	submit.Expiration = c.PostForm("expire")
	submit.Password = c.PostForm("password")
	submit.uploadPassword = c.PostForm("upload_password")
	submit.Syntax = c.PostForm("syntax")
	submit.Privacy = c.PostForm("privacy")
	burnAfterInt, err := strconv.Atoi(c.PostForm("burnafter"))
	if err != nil {
		return Submit{}, errors.New("burnafter must be an integer")
	}
	submit.BurnAfter = burnAfterInt

	form, err := c.MultipartForm()
	if err != nil {
		return Submit{}, errors.New("form error: " + err.Error())
	}

	submit.Files = form.File["files[]"]
	return submit, nil
}

func validateSubmit(submit *Submit) error {
	hasFiles := len(submit.Files) > 0
	if submit.Text == "" && !hasFiles {
		return errors.New("text or files is required")
	}

	if 0 < len(submit.Files) && !config.Vars.FileUpload {
		return errors.New("file uploads are disabled")
	}

	if config.Vars.FileUploadingPassword != "" && submit.uploadPassword == "" || (submit.uploadPassword != config.Vars.FileUploadingPassword) {
		return errors.New("invalid upload password")
	}

	needsAuth := (submit.Privacy == "private" || submit.Privacy == "secret" || submit.Privacy == "readonly")
	if submit.Password == "" && needsAuth {
		return errors.New("password is required for private, secret or readonly pastes")
	}

	if 128 < len(submit.Password) {
		return errors.New("keep them passwords under sane lengths :)")
	}

	if !paste.PrivacyMap.Contains(submit.Privacy) {
		return errors.New("invalid privacy")
	}

	if !paste.SyntaxMap.Contains(submit.Syntax) {
		return errors.New("invalid syntax")
	}

	if !config.Vars.EternalPaste && submit.Expiration == "never" {
		submit.Expiration = "1w"
	}

	if config.Vars.MaxContentLength < len(submit.Text) {
		return errors.New("content is too long")
	}

	maxSizeFiles := utils.Ternary((needsAuth && submit.Privacy != "readonly"), config.Vars.MaxEncryptionSize, config.Vars.MaxFileSize)

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
