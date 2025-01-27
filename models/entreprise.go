package models

import (
	"time"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type Entreprise struct {
	gorm.Model

	TypeEntreprise string    `gorm:"not null" json:"type_entreprise"` // PME, GE, Particulier
	Name           string    `gorm:"not null" json:"name"`
	Code           string    `gorm:"not null" json:"code"` // Code entreprise
	Rccm           string    `json:"rccm"`
	IdNat          string    `json:"idnat"`
	NImpot         string    `json:"nimpot"`
	Adresse        string    `json:"adresse"`
	Email          string    `json:"email"`                     // Email officiel
	Telephone      string    `gorm:"not null" json:"telephone"` // Telephone officiel
	Manager        string    `gorm:"not null" json:"manager"`
	Status         bool      `gorm:"not null" json:"status"`
	TypeAbonnement string    `json:"type_abonnement"`
	Abonnement     time.Time `json:"abonnement"`
	Signature      string    `json:"signature"`
	Users          []User    `gorm:"foreignKey:EntrepriseID"`
	Pos            []Pos     `gorm:"foreignKey:EntrepriseID"`
}

type EntrepriseInfos struct {
	ID             uint      `json:"id"`
	TypeEntreprise string    `json:"type_entreprise"` // PME, GE, Particulier
	Name           string    `json:"name"`
	Code           string    `json:"code"` // Code entreprise
	Rccm           string    `json:"rccm"`
	IdNat          string    `json:"idnat"`
	NImpot         string    `json:"nimpot"`
	Adresse        string    `json:"adresse"`
	Email          string    `json:"email"`     // Email officiel
	Telephone      string    `json:"telephone"` // Telephone officiel
	Manager        string    `json:"manager"`
	Status         bool      `json:"status"`
	TypeAbonnement string    `json:"type_abonnement"`
	Abonnement     time.Time `json:"abonnement"`
	Signature      string    `json:"signature"`
	TotalUser      int       `json:"total_user"`
	TotalPos       int       `json:"total_pos"`
}
