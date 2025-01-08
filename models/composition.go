package models

import "gorm.io/gorm"

type Composition struct {
	gorm.Model

	PosID          uint       `json:"pos_id"`
	PlatID         uint       `json:"plat_id"`
	IngredientID   uint       `json:"ingredient_id"`
	Quantity       uint64     `gorm:"not null" json:"quantity"`
	Signature      string     `json:"signature"`
	CodeEntreprise uint64     `json:"code_entreprise"`
	Plat           Plat       `gorm:"foreignKey:PlatID"`
	Ingredient     Ingredient `gorm:"foreignKey:IngredientID"`
	Pos            Pos        `gorm:"foreignKey:PosID"`
}
