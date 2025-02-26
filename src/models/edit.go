package models

type Edit struct {
	CommonFields
	Type            string `json:"type"`
	Ts              int64  `json:"ts"`
	RemoteFileModTs string `json:"remoteFileModTs"`
	Key             string `json:"key"`
	Value           string `json:"value"`                       // For added keys
	NewValue        string `json:"newValue"`                    // For existing keys
	OldValue        string `json:"oldValue"`                    // For existing keys
	JobID           uint   `gorm:"not null;index" json:"jobId"` // Edit N:1 Job
}
