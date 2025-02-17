package models

import "gorm.io/gorm"

type Stat struct {
	gorm.Model
	JsonStr string `gorm:"type:json"`
}
