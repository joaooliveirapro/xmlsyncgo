package models

type Job struct {
	CommonFields
	ExternalReference string `json:"externalReference"`
	Deleted           bool   `json:"-"`
	Content           string `gorm:"type:json" json:"content"`
	Hash              string `json:"hash"`
	FileID            uint   `gorm:"not null;index" json:"fileId"` // Job N:1 File
	Edits             []Edit `gorm:"foreignKey:JobID" json:"-"`    // Job 1:N Edit
}
