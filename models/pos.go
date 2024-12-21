package models

import (
	"gorm.io/gorm"
)

type Pos struct {
	gorm.Model

	EntrepriseID uint          `json:"entreprise_id"`
	Entreprise   Entreprise    `gorm:"foreignKey:EntrepriseID"`
	Name         string        `gorm:"not null" json:"name"`
	Adresse      string        `json:"adresse"`
	Email        string        `json:"email"`
	Telephone    string        `json:"telephone"`
	Manager      string        `gorm:"not null" json:"manager"`
	Status       bool          `gorm:"not null" json:"status"` // Actif ou Inactif
	Signature    string        `json:"signature"`
	Stocks       []Stock       `gorm:"foreignKey:PosID" json:"stocks"` 
	Commandes    []Commande    `gorm:"foreignKey:PosID" json:"commandes"`
}
