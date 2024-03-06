package main

type SubmittedPassword struct {
	Password string `form:"password"`
}

type PasteLists struct {
	Pastes    []Paste
	Redirects []Paste
}

type File struct {
	ID   int
	Name string
	Size int
	Blob []byte
}

type Paste struct {
	ID             int
	PasteName      string
	Expire         string
	Privacy        string
	ReadCount      int
	ReadLast       string
	BurnAfter      int
	Content        string
	UrlRedirect    int
	Syntax         string
	HashedPassword string
	Files          []File
	CreatedAt      string
	UpdatedAt      string
}

type ConfigEnv struct {
	Salt          string
	Port          string
	DataDir       string
	AdminPassword string

	PublicList bool
	PublicURL  string

	DefaultExpiryTime   string
	NoFileUpload        bool
	MaxFileSize         int
	MaxEncryptionSize   int
	MaxContentLength    int
	UploadingPassword   string
	DisableEternalPaste bool
	DisableReadCount    bool
	DisableBurnAfter    bool
	DefaultExpiry       string
	ShortPasteNames     bool
	IUnderstandTheRisks bool
}
