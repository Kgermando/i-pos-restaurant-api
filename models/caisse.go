package models

import "gorm.io/gorm"

type Caisse struct {
	gorm.Model

	TypeTransaction string  `gorm:"not null" json:"type_transaction"` // Entrée ou Sortie
	Montant         float64 `gorm:"not null" json:"montant"`          // Montant de la transaction
	Libelle         string  `json:"libelle"`          // Description de la transaction
	Reference       string  `json:"reference"`        // Nombre aleatoire
	Signature       string  `json:"signature"`        // Signature de la transaction
	PosID           uint    `json:"pos_id"`           // ID du point de vente
	Pos             Pos     `gorm:"foreignKey:PosID"` // Point de vente
	CodeEntreprise  uint    `json:"code_entreprise"`  // ID de l'entreprise
}
