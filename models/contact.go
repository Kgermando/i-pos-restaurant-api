package models

import (
	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model

	// ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Fullname string    `gorm:"not null" json:"fullname"`
	Email    string    `gorm:"not null" json:"email"`
	Subject  string    `gorm:"not null" json:"subject"`
	Message  string    `gorm:"not null" json:"message"`
}
