package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User ...
type User struct {
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	ID        string     `gorm:"primary_key"`
	UpdatedAt time.Time
	Name      string
	Password  string
	Email     string
	Phone     string
	Friends   []User
}

// Friend ... this will keep track of id that who is connected with whom
type Friendship struct {
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	ID        string     `gorm:"primary_key"`
	ToUser    string
	FromUser  string
	Status    string
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
