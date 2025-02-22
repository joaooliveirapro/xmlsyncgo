package models

import (
	"time"
)

type File struct {
	CommonFields
	RemoteModTime        time.Time  `json:"remoteModTim"`
	RootKey              string     `json:"rootKey"`
	JobNodeKey           string     `json:"jobNodeKey"`
	ExternalReferenceKey string     `json:"externalReferenceKey"`
	RemoteFilepath       string     `json:"remoteFilepath"`
	RemoteFilename       string     `json:"remoteFilename"`
	Hostname             string     `json:"hostname"`
	Port                 string     `json:"port"`
	Username             string     `json:"username"`
	Password             string     `json:"password"`
	AuditIteration       uint       `json:"auditIteration"`
	ClientID             uint       `gorm:"not null;index" json:"clientId"` // File N:1 Client
	Audit                []AuditLog `gorm:"foreignKey:FileID" json:"audit"` // File 1:N AuditLog
	Stats                []Stat     `gorm:"foreignKey:FileID" json:"stats"` // File 1:N Stat
	Jobs                 []Job      `gorm:"foreignKey:FileID" json:"jobs"`  // File 1:N Job
}
