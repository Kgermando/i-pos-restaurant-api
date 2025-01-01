package models

import "gorm.io/gorm"

type Area struct {
	gorm.Model

	Name           string      `gorm:"not null" json:"name"`
	Province       string      `gorm:"not null" json:"province"`
	Signature      string      `json:"signature"`
	CodeEntreprise uint        `json:"code_entreprise"`
	Livraisons     []Livraison `gorm:"foreignKey:AreaID"`
}
