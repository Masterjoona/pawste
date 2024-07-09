package utils

import (
	"mime/multipart"

	"github.com/Masterjoona/pawste/pkg/paste"
)

type Submit struct {
	Text       string                  `form:"text,omitempty"`
	Expiration string                  `form:"expiration,omitempty"`
	BurnAfter  int                     `form:"burn,omitempty"`
	Password   string                  `form:"password,omitempty"`
	Syntax     string                  `form:"syntax,omitempty"`
	Privacy    string                  `form:"privacy,omitempty"`
	Files      []*multipart.FileHeader `form:"file,omitempty"`
}

type PasteUpdate struct {
	Content        string                  `form:"content,omitempty"`
	Password       string                  `form:"password,omitempty"`
	RemovedFiles   []string                `form:"removed_files,omitempty"`
	FilesMultiPart []*multipart.FileHeader `form:"file,omitempty"`
	Files          []paste.File
}
