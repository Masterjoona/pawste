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
	Expire      int64
	Privacy     string
	IsEncrypted int
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
