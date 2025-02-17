package models

import (
	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	Text   string
	FileID uint
}

func NewAuditLog(text string) error {
	result := initializers.DB.Create(&AuditLog{Text: text})
	return result.Error
}
