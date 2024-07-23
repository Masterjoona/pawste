package paste

import "mime/multipart"

type File struct {
	ID          int
	Name        string
	Size        int
	ContentType string
	Blob        []byte
}

type Paste struct {
	ID          int
	PasteName   string
	Expire      int64
	Privacy     string
	NeedsAuth   int
	ReadCount   int
	ReadLast    int64
	BurnAfter   int
	Content     string
	UrlRedirect int
	Syntax      string
	Password    string
	Files       []File
	CreatedAt   int64
	UpdatedAt   int64
}

type PasteLists struct {
	Pastes    []Paste
	Redirects []Paste
}

type PasteUpdate struct {
	Content        string                  `form:"content,omitempty"`
	Password       string                  `form:"password,omitempty"`
	RemovedFiles   []string                `form:"removed_files,omitempty"`
	FilesMultiPart []*multipart.FileHeader `form:"file,omitempty"`
	Files          []File
}

var PrivacyOptions = []string{"public", "unlisted", "readonly", "private", "secret"}
