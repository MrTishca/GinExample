package core

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string `json:"name"`
	Email         string `gorm:"unique" json:"email"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	LastConection string `json:"last_conection"`
	TimeSession   string `json:"time_session"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}
