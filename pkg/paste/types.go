package paste

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
	Expire      string
	Privacy     string
	IsEncrypted int
	ReadCount   int
	ReadLast    string
	BurnAfter   int
	Content     string
	UrlRedirect int
	Syntax      string
	Password    string
	Files       []File
	CreatedAt   string
	UpdatedAt   string
}

type PasteLists struct {
	Pastes    []Paste
	Redirects []Paste
}
