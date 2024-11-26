package model

import (
	"gorm.io/gorm"
)

// RevokedToken represents a revoked token in the database
type RevokedToken struct {
	gorm.Model
	Token string `gorm:"not null;uniqueIndex"` // Token yang dibatalkan
}
