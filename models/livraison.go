package models

import "gorm.io/gorm"

type Livraison struct {
	gorm.Model

	OperatorName   string  `gorm:"not null" json:"operator_name"`
	AreaID         uint    `json:"area_id"`
	CoutLivraison  float64 `gorm:"not null" json:"cout_livraison"`
	ClientID       uint    `json:"client_id"`
	LivreurID      uint    `json:"livreur_id"`
	PosID          uint    `json:"pos_id"`
	Status         string  `json:"status"`
	Signature      string  `json:"signature"`
	CodeEntreprise uint64  `json:"code_entreprise"`

	Client  Client  `gorm:"foreignKey:ClientID"`
	Livreur Livreur `gorm:"foreignKey:LivreurID"`
	Area    Area    `gorm:"foreignkey:AreaID"`
	Pos     Pos     `gorm:"foreignKey:PosID"`

	CommandeLines []CommandeLine `gorm:"foreignKey:LivraisonID"`
}

type LivraisonArea struct {
	Fullname  string `json:"fullname"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Count     string `json:"count"`
}
