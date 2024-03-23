package main

import (
	"io"
	"mime/multipart"
	"os"
	"regexp"
	"time"
)

func SubmitToPaste(submit Submit, pasteName string, isRedirect int) Paste {
	var files []File
	for _, file := range submit.Files {
		if file == nil {
			continue
		}
		fileName, fileSize, fileBlob := convertMultipartFile(file)
		files = append(files, File{
			Name: fileName,
			Size: fileSize,
			Blob: fileBlob,
		})
	}

	return Paste{
		PasteName:      pasteName,
		Expire:         HumanTimeToSQLTime(submit.Expiration),
		Privacy:        submit.Privacy,
		IsEncrypted:    0,
		ReadCount:      1,
		ReadLast:       GetCurrentDate(),
		BurnAfter:      submit.BurnAfter,
		Content:        submit.Text,
		Syntax:         submit.Syntax,
		HashedPassword: submit.Password,
		Files:          files,
		UrlRedirect:    isRedirect,
		CreatedAt:      GetCurrentDate(),
		UpdatedAt:      GetCurrentDate(),
	}
}

func convertMultipartFile(file *multipart.FileHeader) (string, int, []byte) {
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

func HumanTimeToSQLTime(humanTime string) string {
	var duration time.Duration
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
	case "1w":
		duration = 7 * 24 * time.Hour
	case "never":
		duration = 100 * 365 * 25 // cope if you're still using this in 100 years
	default:
		duration = 7 * 24 * time.Hour
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

func SaveFileToDisk(file *File, pasteName string) error {
	err := os.WriteFile(
		Config.DataDir+pasteName+"/"+file.Name,
		file.Blob,
		0644,
	)
	if err != nil {
		return err
	}
	return nil
}
