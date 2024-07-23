package handling

import (
	"mime/multipart"
)

type Submit struct {
	Text           string
	Expiration     string
	BurnAfter      int
	Password       string
	uploadPassword string
	Syntax         string
	Privacy        string
	Files          []*multipart.FileHeader
}
