package models

import (
	"time"

	"gorm.io/gorm"
)

type IngredientStock struct {
	gorm.Model

	PosID          uint        `json:"pos_id"`
	Pos            Pos         `gorm:"foreignKey:PosID"`
	IngredientID   uint        `json:"ingredient_id"`
	Ingredient     Ingredient  `gorm:"foreignKey:IngredientID"`
	Description    string      `json:"description"`
	Quantity       uint64      `gorm:"not null" json:"quantity"`
	FournisseurID  uint        `json:"fournisseur_id"`
	Fournisseur    Fournisseur `gorm:"foreignKey:FournisseurID"`
	PrixAchat      float64     `gorm:"not null" json:"prix_achat"`
	DateExpiration time.Time   `gorm:"not null" json:"date_expiration"`
	Signature      string      `json:"signature"`
	CodeEntreprise uint64      `json:"code_entreprise"`
}

type IngredientUsage struct {
	Id             uint    `gorm:"column:id"`
	PlatID         uint    `gorm:"column:plat_id"`
	IngredientName string  `gorm:"column:name"`
	Qty            float64 `gorm:"column:qty"`
	Unite          string  `gorm:"column:unite"`
	Price          float64 `gorm:"column:price"`
}
