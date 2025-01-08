package models

import "gorm.io/gorm"

type Livreur struct {
	gorm.Model

	NameSociety    string `gorm:"not null" json:"name_society"`
	LivreurName    string `gorm:"not null" json:"livreur_name"`
	Telephone      string `gorm:"not null" json:"telephone"`
	Email          string `json:"email"`
	Rccm           string `json:"rccm"`
	IdNat          string `json:"idnat"`
	Signature      string `json:"signature"`
	CodeEntreprise uint64 `json:"code_entreprise"`

	Livraison []Livraison `gorm:"foreignKey:LivreurID"`
}
