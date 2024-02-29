package main

import (
	"io"
	"mime/multipart"
	"os"
	"regexp"
	"time"
)

func SubmitToPaste(submit Submit, pasteName string, hashedPassword string) Paste {
	var files []File
	for _, file := range submit.Files {
		if file == nil {
			continue
		}
		fileName, fileSize, fileBlob := multipartIntoThings(file)
		files = append(files, File{
			FileName: fileName,
			FileSize: fileSize,
			FileBlob: fileBlob,
		})
	}
	return Paste{
		PasteName:      pasteName,
		Expire:         HumanTimeToSQLTime(submit.Expiration),
		Privacy:        submit.Privacy,
		ReadCount:      1,
		ReadLast:       GetCurrentDate(),
		BurnAfter:      submit.BurnAfter,
		Content:        submit.Text,
		Syntax:         submit.Syntax,
		HashedPassword: hashedPassword,
		Files:          files,
		UrlRedirect:    IsContentJustUrl(submit.Text),
		CreatedAt:      GetCurrentDate(),
		UpdatedAt:      GetCurrentDate(),
	}
}

func multipartIntoThings(file *multipart.FileHeader) (string, int, []byte) {
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
	// time into 2024-02-26T15:56:16Z + duration
	return time.Now().Add(duration).Format("2006-01-02 15:04:05")
}
func DoesPasteHaveFiles(paste Paste) bool {
	if _, err := os.Stat(Config.DataDir + paste.PasteName); os.IsNotExist(err) {
		return false
	}

	files, err := os.ReadDir(Config.DataDir + paste.PasteName)
	if err != nil {
		return false
	}
	return len(files) > 0
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
