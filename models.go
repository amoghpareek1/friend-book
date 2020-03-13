package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User ...
type User struct {
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	ID        string     `gorm:"type:uuid;primary_key;"`
	UpdatedAt time.Time
	FirstName string
	LastName  string
	Password  string
	Email     string
}

// Response ...
type Response struct {
	Success bool
	Data    interface{}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
