package models

import (
	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model

	// ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Fullname       string     `gorm:"not null" json:"fullname"`
	Telephone      string     `gorm:"not null" json:"telephone"`
	Email          string     `json:"email"`
	Signature      string     `json:"signature"`
	Commandes      []Commande `gorm:"foreignKey:ClientID"`
	CodeEntreprise uint       `json:"code_entreprise"`
}
