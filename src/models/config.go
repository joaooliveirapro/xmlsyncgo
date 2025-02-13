package models

type ConfigModel struct {
	Root                 string          `json:"root"`
	JobNode              string          `json:"job_node"`
	ExternalReferenceKey string          `json:"external_reference_key"`
	RemoteFilepath       string          `json:"remote_filepath"`
	RemoteFilename       string          `json:"remote_filename"`
	LocalFileTempFolder  string          `json:"download_to"`
	SftpCredentials      SftpCredentials `json:"sftp_credentials"`
}

type SftpCredentials struct {
	Hostname string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
