package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model

	Fullname       string `gorm:"not null" json:"fullname"`
	Telephone      string `gorm:"not null" json:"telephone"`
	Email          string `json:"email"`
	Adress         string `json:"adress"`
	Signature      string `json:"signature"`
	CodeEntreprise uint   `json:"code_entreprise"`
	

	Livraison []Livraison `gorm:"foreignKey:ClientID"`
	Commandes []Commande  `gorm:"foreignKey:ClientID"`
}
 