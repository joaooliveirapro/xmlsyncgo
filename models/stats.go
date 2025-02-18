package models

import "gorm.io/gorm"

type Stat struct {
	gorm.Model
	JsonStr string `gorm:"type:json"`
	FileID  uint   `gorm:"not null;index"` // Stat N:1 File
}
