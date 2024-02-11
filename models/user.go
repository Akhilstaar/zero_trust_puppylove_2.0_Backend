package models

import (
	"gorm.io/gorm"
)

type (
	// User represents the structure of our resource
	User struct {
		gorm.Model
		Id        string `json:"_id" bson:"_id" gorm:"unique"`
		Name      string `json:"name" bson:"name"`
		Email     string `json:"email" bson:"email" gorm:"unique"`
		Pass      string `json:"passHash" bson:"passHash"`
		PubK      string `json:"pubKey" bson:"pubKey"`
		PrivK     string `json:"privKey" bson:"privKey"`
		AuthC     string `json:"authCode" bson:"authCode"`
		Data      string `json:"data" bson:"data"`
		S1submit  bool   `json:"s1submit" bson:"s1submit"`
		S2submit  bool   `json:"s2submit" bson:"s2submit"`
		Dirty     bool   `json:"dirty" bson:"dirty"`
		Certgiven bool   `json:"certgiven" bson:"certgiven"`
	}
)

type UserPublicKey struct {
	Id   string `json:"_id" bson:"_id" gorm:"unique"`
	PubK string `json:"pubKey" bson:"pubKey"`
}

type UserLogin struct {
	Id   string `json:"id" binding:"required"`
	Pass string `json:"passHash" binding:"required"`
}

type AddNewUser struct {
	TypeUserNew []TypeUserNew `json:"newuser" binding:"required"`
}

type TypeUserNew struct {
	Id    string `json:"roll" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type TypeUserFirst struct {
	Id       string `json:"roll" binding:"required"`
	AuthCode string `json:"authCode" binding:"required"`
	PassHash string `json:"passHash" binding:"required"`
	PubKey   string `json:"pubKey" binding:"required"`
	PrivKey  string `json:"privKey" binding:"required"`
	Data     string `json:"data" binding:"required"` // TODO: Add limit to it's size, else someone will fill it with really large junk data & slow down DB
}

// w'll change it later (maybee..)
type AdminLogin struct {
	Id   string `json:"id" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}

type MailData struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	AuthC string `json:"authCode" binding:"required"`
	Dirty bool   `json:"dirty" binding:"required"`
}

var RegisterMap = make(map[string]int)
