package config

type ConfigEnv struct {
	Salt          string
	Port          string
	DataDir       string
	AdminPassword string

	PublicList bool
	PublicURL  string

	DefaultExpiryTime     string
	FileUpload            bool
	MaxFileSize           int
	MaxEncryptionSize     int
	MaxContentLength      int
	UploadingPassword     string
	EternalPaste          bool
	ReadCount             bool
	BurnAfter             bool
	DefaultExpiry         string
	ShortPasteNames       bool
	ShortenRedirectPastes bool
	CountFileUsage        bool
}
