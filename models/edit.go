package models

import "gorm.io/gorm"

type Edit struct {
	gorm.Model
	Type            string
	Ts              int64
	RemoteFileModTs string
	Key             string
	Value           string // For added keys
	NewValue        string // For existing keys
	OldValue        string // For existing keys
	JobID           uint   `gorm:"not null;index"` // Edit N:1 Job
}
