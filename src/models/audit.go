package models

type AuditLog struct {
	CommonFields
	Text           string `json:"text"`
	FileID         uint   `gorm:"not null;index" json:"fileId"` // AuditLog N:1 File
	AuditIteration uint   `json:"auditIteration"`
}
