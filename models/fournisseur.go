package models

import (
	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type Fournisseur struct {
	gorm.Model

	// ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name           string    `gorm:"not null" json:"name"` // Name of entreprise
	Telephone      string    `json:"telephone"`
	Email          string    `json:"email"`
	Adresse        string    `json:"adresse"`
	TypeFourniture string    `json:"type_fourniture"`
	Signature      string    `json:"signature"`
	CodeEntreprise uint      `json:"code_entreprise"`
	Stocks         []Stock   `gorm:"foreignKey:FournisseurID"`
}
