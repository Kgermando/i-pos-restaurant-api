package models

import (
	"gorm.io/gorm"
)

type CommandeLine struct {
	gorm.Model

	CommandeID  uint      `json:"commande_id"`
	Commande    Commande  `gorm:"foreignKey:CommandeID"`
	LivraisonID uint      `json:"livraison_id"` // LivraisonID est identique a commande dans le fonctionnement
	Livraison   Livraison `gorm:"foreignKey:LivraisonID"`

	ProductID      uint    `json:"product_id"`
	PlatID         uint    `json:"plat_id"`
	Product        Product `gorm:"foreignKey:ProductID"`
	Plat           Plat    `gorm:"foreignKey:PlatID"`
	Quantity       uint64  `gorm:"not null" json:"quantity"`
	CodeEntreprise uint64  `json:"code_entreprise"`
}
