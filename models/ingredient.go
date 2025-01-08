package models

import "gorm.io/gorm"

type Ingredient struct {
	gorm.Model

	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
	// PrixUnitaire   float64       `gorm:"not null" json:"prix_unitaire"`
	Unite          string        `gorm:"not null" json:"unite"`
	PosID          uint          `json:"pos_id"`
	Pos            Pos           `gorm:"foreignKey:PosID"`
	Signature      string        `json:"signature"`
	CodeEntreprise uint64        `json:"code_entreprise"`
	Compositions   []Composition `gorm:"foreignKey:IngredientID"`
}
