package models

import (
	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	Text           string
	FileID         uint `gorm:"not null;index"` // AuditLog N:1 File
	AuditIteration uint
}
