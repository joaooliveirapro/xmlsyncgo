package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model // Already includes some fields (ID, Created, Updated, Deleted ts)
	Name       string
	Files      []File // File -*--1- Client
}
