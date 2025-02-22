package models

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Stat struct {
	CommonFields
	JsonStr string         `gorm:"type:json" json:"jsonStr"`
	FileID  uint           `gorm:"not null;index" json:"fileId"` // Stat N:1 File
	content map[string]int // Not exported (lowercase)
}

func NewStat(keys []string, fileID uint) *Stat {
	s := Stat{FileID: fileID, content: map[string]int{}}
	for _, key := range keys {
		s.content[key] = 0 // Set keys to 0
	}
	return &s
}

func (s *Stat) AddEntry(key string, val int) {
	s.content[key] = val
}

func (s *Stat) IncrementKey(key string, by int) {
	s.content[key] += by
}

func (s *Stat) SaveToDB(tx *gorm.DB) error {
	// Convert stat map to json string
	statJson, err := json.Marshal(s.content)
	if err != nil {
		return fmt.Errorf("[stats.go:SaveToDB] error converting statJson to JSON. %w", err)
	}
	// Assign the statJson string to the model
	s.JsonStr = string(statJson)
	// Save model to DB
	result := tx.Save(&s)
	if result.Error != nil {
		return fmt.Errorf("[stats.go:SaveToDB] error saving stat %w", result.Error)
	}
	return nil
}
