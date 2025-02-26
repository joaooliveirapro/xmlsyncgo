package models

type Client struct {
	CommonFields
	Name  string `json:"name"`
	Files []File `gorm:"foreignKey:ClientID" json:"-"`
}
