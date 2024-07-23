package utils

import (
	"io"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/romana/rlog"
)

func SubmitToPaste(submit Submit, pasteName string, isRedirect int) (paste.Paste, error) {
	var files []paste.File
	for _, file := range submit.Files {
		if file == nil {
			continue
		}
		fileName, fileSize, fileBlob, err := ConvertMultipartFile(file)
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
		Expire:      humanTimeToUnix(submit.Expiration),
		Privacy:     submit.Privacy,
		NeedsAuth:   Ternary((submit.Password == ""), 0, 1).(int),
		ReadCount:   0,
		ReadLast:    todaysDate,
		BurnAfter:   Ternary(config.Vars.BurnAfter, submit.BurnAfter, 0).(int),
		Content:     submit.Text,
		Syntax:      submit.Syntax,
		Password:    submit.Password,
		Files:       files,
		UrlRedirect: isRedirect,
		CreatedAt:   todaysDate,
		UpdatedAt:   todaysDate,
	}, nil
}

func ConvertMultipartFile(file *multipart.FileHeader) (string, int, []byte, error) {
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

func humanTimeToUnix(humanTime string) int64 {
	if humanTime == "never" {
		return -1
	}
	duration := config.ParseDuration(humanTime)
	if config.ParseDuration(config.Vars.MaxExpiryTime) < duration {
		return time.Now().Add(time.Duration(config.OneWeek)).Unix()
	}
	return time.Now().Add(duration).Unix()
}

func IsContentJustUrl(content string) int {
	if regexp.MustCompile(`^(?:http|https|magnet):\/\/[^\s/$.?#].[^\s]*$`).MatchString(content) {
		return 1
	}
	return 0
}

func AllowedOption(s string, options []string) bool {
	for _, item := range options {
		if item == s {
			return true
		}
	}
	return false
}

func Ternary(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}
