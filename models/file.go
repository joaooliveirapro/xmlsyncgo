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
	ClientID             uint // File -*--1- Client
	AuditIteration       uint
	Audit                []AuditLog
}
