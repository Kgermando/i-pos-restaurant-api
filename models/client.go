package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model

	Fullname     string `gorm:"not null" json:"fullname"`
	Telephone    string `gorm:"not null" json:"telephone"`
	Telephone2   string `json:"telephone2"`
	Email        string `json:"email"`
	Adress       string `json:"adress"`
	// Birthday     string `json:"birthday"`
	Organisation string `json:"organisation"`
	WebSite      string `json:"website"`

	Signature      string `json:"signature"`
	CodeEntreprise uint64 `json:"code_entreprise"`

	Livraison []Livraison `gorm:"foreignKey:ClientID"`
	Commandes []Commande  `gorm:"foreignKey:ClientID"`
}
