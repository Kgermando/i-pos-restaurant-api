package models

import (
	"gorm.io/gorm"
)

type Plat struct {
	gorm.Model

	Reference      string         `gorm:"not null" json:"reference"`
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `gorm:"not null" json:"description"`
	UniteVente     string         `json:"unite_vente"`
	PrixVente      float64        `gorm:"not null" json:"prix_vente"`
	Tva            float64        `gorm:"default:0" json:"tva"`
	Signature      string         `json:"signature"`
	PosID          uint           `json:"pos_id"`
	Pos            Pos            `gorm:"foreignKey:PosID"`
	CodeEntreprise uint           `json:"code_entreprise"`
	CommadeLines   []CommandeLine `gorm:"foreignKey:ProductID"`
}
