package models

type ConfigModel struct {
	Id                   int64
	Root                 string
	JobNode              string
	ExternalReferenceKey string
	RemoteFilepath       string
	RemoteFilename       string
	LocalFileTempFolder  string
	SftpCredentials      SftpCredentials
}

type SftpCredentials struct {
	Hostname string
	Port     string
	Username string
	Password string
}
