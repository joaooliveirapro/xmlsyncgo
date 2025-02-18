package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	RemoteModTime        time.Time
	RootKey              string
	JobNodeKey           string
	ExternalReferenceKey string
	RemoteFilepath       string
	RemoteFilename       string
	Hostname             string
	Port                 string
	Username             string
	Password             string
	AuditIteration       uint
	ClientID             uint       `gorm:"not null;index"`    // File N:1 Client
	Audit                []AuditLog `gorm:"foreignKey:FileID"` // File 1:N AuditLog
	Stats                []Stat     `gorm:"foreignKey:FileID"` // File 1:N Stat
	Jobs                 []Job      `gorm:"foreignKey:FileID"` // File 1:N Job
}
