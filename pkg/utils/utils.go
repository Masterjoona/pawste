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
	todaysDate := GetCurrentDate()
	return paste.Paste{
		PasteName:   pasteName,
		Expire:      humanTimeToSQLTime(submit.Expiration),
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

func humanTimeToSQLTime(humanTime string) string {
	var duration time.Duration
	duration = 7 * 24 * time.Hour
	switch humanTime {
	case "10min":
		duration = 10 * time.Minute
	case "1min":
		duration = 1 * time.Minute
	case "1h":
		duration = 1 * time.Hour
	case "6h":
		duration = 6 * time.Hour
	case "24h":
		duration = 24 * time.Hour
	case "72h":
		duration = 72 * time.Hour
	case "never":
		duration = 100 * 365 * 24 * time.Hour // cope if you're still using this in 100 years
	}

	return time.Now().Add(duration).Format("2006-01-02 15:04:05")
}

func IsContentJustUrl(content string) int {
	if regexp.MustCompile(`^(?:http|https|magnet):\/\/[^\s/$.?#].[^\s]*$`).MatchString(content) {
		return 1
	}
	return 0
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NotAllowedPrivacy(x string) bool {
	for _, item := range []string{"public", "unlisted", "readonly", "private", "secret"} {
		if item == x {
			return false
		}
	}
	return true
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
