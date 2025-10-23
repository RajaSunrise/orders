package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	ID			uuid.UUID		`json:"id" gorm:"primary_key"`
	Username	string			`json:"username" gorm:"unique"`
	Email		string			`json:"email" gorm:"unique"`
	Password 	string			`json:"password"`
	CreatedAt	time.Time		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt	time.Time		`json:"updated_at" gorm:"autoCreateTime"`
	
}