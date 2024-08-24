package database

import (
	"time"

	"github.com/Masterjoona/pawste/pkg/config"
	"github.com/Masterjoona/pawste/pkg/paste"
	"github.com/Masterjoona/pawste/pkg/utils"
)

func SubmitToPaste(submit paste.Submit, isRedirect int) (paste.Paste, error) {
	var files []paste.File
	for _, file := range submit.Files {
		if file == nil {
			continue
		}
		fileName, fileSize, fileBlob, err := paste.ConvertMultipartFile(file)
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
	pasteName := createPasteName(isRedirect)
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
