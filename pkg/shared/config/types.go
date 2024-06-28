package config

type ConfigEnv struct {
	Salt          string
	Port          string
	DataDir       string
	AdminPassword string

	PublicList bool
	PublicURL  string

	DefaultExpiryTime     string
	NoFileUpload          bool
	MaxFileSize           int
	MaxEncryptionSize     int
	MaxContentLength      int
	UploadingPassword     string
	DisableEternalPaste   bool
	DisableReadCount      bool
	DisableBurnAfter      bool
	DefaultExpiry         string
	ShortPasteNames       bool
	ShortenRedirectPastes bool
	IUnderstandTheRisks   bool
}
