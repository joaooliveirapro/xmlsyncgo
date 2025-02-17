package models

import (
	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	Text           string
	FileID         uint
	AuditIteration uint
}

func NewAuditLog(text string, auditI uint) error {
	result := initializers.DB.Create(&AuditLog{Text: text, AuditIteration: auditI})
	return result.Error
}
