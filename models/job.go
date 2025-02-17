package models

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	ExternalReference string
	Content           string `gorm:"type:json"`
	Hash              string
	Edits             []Edit
}
