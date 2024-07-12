package utils

import (
	"io"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
)

func SubmitToPaste(submit Submit, pasteName string, isRedirect int) paste.Paste {
	var files []paste.File
	for _, file := range submit.Files {
		if file == nil {
			continue
		}
		fileName, fileSize, fileBlob := ConvertMultipartFile(file)
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
		IsEncrypted: TernaryInt((submit.Password != ""), 1, 0),
		ReadCount:   0,
		ReadLast:    todaysDate,
		BurnAfter:   TernaryInt(config.Config.BurnAfter, submit.BurnAfter, 0),
		Content:     submit.Text,
		Syntax:      submit.Syntax,
		Password:    submit.Password,
		Files:       files,
		UrlRedirect: isRedirect,
		CreatedAt:   todaysDate,
		UpdatedAt:   todaysDate,
	}
}

func ConvertMultipartFile(file *multipart.FileHeader) (string, int, []byte) {
	src, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()

	fileBlob, err := io.ReadAll(src)
	if err != nil {
		panic(err)
	}
	return file.Filename, len(fileBlob), fileBlob
}

func humanTimeToUnix(humanTime string) int64 {
	duration := time.Duration(config.ParseDuration(humanTime))
	if time.Duration(config.Config.MaxExpiryTime) < duration {
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

func TernaryString(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func TernaryInt(condition bool, trueVal, falseVal int) int {
	if condition {
		return trueVal
	}
	return falseVal
}
