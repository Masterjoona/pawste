package paste

import (
	"io"
	"mime/multipart"

	"github.com/Masterjoona/pawste/pkg/config"
)

func ConvertMultipartFile(file *multipart.FileHeader) (string, int, []byte, error) {
	src, err := file.Open()
	if err != nil {
		config.Logger.Error("Could not open multipart file", err)
		return "", 0, nil, err
	}
	defer src.Close()

	fileBlob, err := io.ReadAll(src)
	if err != nil {
		config.Logger.Error("Could not read multipart file", err)
		return "", 0, nil, err
	}
	return file.Filename, len(fileBlob), fileBlob, nil
}
