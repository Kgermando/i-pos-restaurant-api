package models

import (
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model

	PosID          uint        `json:"pos_id"`
	Pos            Pos         `gorm:"foreignKey:PosID"`
	ProductID      uint        `json:"product_id"`
	Product        Product     `gorm:"foreignKey:ProductID"`
	Description    string      `json:"description"`
	Quantity       uint64      `gorm:"not null" json:"quantity"`
	FournisseurID  uint        `json:"fournisseur_id"`
	Fournisseur    Fournisseur `gorm:"foreignKey:FournisseurID"`
	PrixAchat      float64     `gorm:"not null" json:"prix_achat"`
	DateExpiration time.Time   `gorm:"not null" json:"date_expiration"`
	Signature      string      `json:"signature"`
	CodeEntreprise uint64      `json:"code_entreprise"`
}

type FournisseurStock struct {
	Name           string  `json:"name"`
	Telephone      string  `json:"telephone"`
	TypeFourniture string  `json:"type_fourniture"`
	TotalValue     float64 `json:"total_value"`
}
