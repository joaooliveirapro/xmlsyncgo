package models

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	ExternalReference string
	Deleted           bool
	Content           string `gorm:"type:json"`
	Hash              string
	FileID            uint   `gorm:"not null;index"`   // Job N:1 File
	Edits             []Edit `gorm:"foreignKey:JobID"` // Job 1:N Edit

}
