package models

import (
	"gorm.io/gorm"
)

type TableBox struct {
	gorm.Model

	PosID          uint       `json:"pos_id"`
	Pos            Pos        `gorm:"foreignKey:PosID"`
	Name           string     `gorm:"not null" json:"name"`
	Numero         int        `gorm:"not null" json:"numero"`
	Status         string     `json:"status"` // Ouverte et Fermée
	Signature      string     `json:"signature"`
	CodeEntreprise uint       `json:"code_entreprise"`
	Commandes      []Commande `gorm:"foreignKey:TableBoxID"`
}
