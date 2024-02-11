package models

import (
	"gorm.io/gorm"
)

type (
	Stage1 struct {
		gorm.Model
		Id string `json:"id" bson:"id" gorm:"unique"`
		M1 string `json:"m1" bson:"m1" gorm:"unique"`
		M2 string `json:"m2" bson:"m2" gorm:"unique"`
		M3 string `json:"m3" bson:"m3" gorm:"unique"`
		M4 string `json:"m4" bson:"m4" gorm:"unique"`
	}
)

type (
	Stage2 struct {
		gorm.Model
		Id    string `json:"id" bson:"id" gorm:"unique"`
		Ka1   string `json:"ka1" bson:"ka1"`
		Ka2   string `json:"ka2" bson:"ka2"`
		Ka3   string `json:"ka3" bson:"ka3"`
		Ka4   string `json:"ka4" bson:"ka4"`
		Kb1   string `json:"kb1" bson:"kb1"`
		Kb2   string `json:"kb2" bson:"kb2"`
		Kb3   string `json:"kb3" bson:"kb3"`
		Kb4   string `json:"kb4" bson:"kb4"`
		Kc1   string `json:"kc1" bson:"kc1"`
		Kc2   string `json:"kc2" bson:"kc2"`
		Kc3   string `json:"kc3" bson:"kc3"`
		Kc4   string `json:"kc4" bson:"kc4"`
		Kd1   string `json:"kd1" bson:"kd1"`
		Kd2   string `json:"kd2" bson:"kd2"`
		Kd3   string `json:"kd3" bson:"kd3"`
		Kd4   string `json:"kd4" bson:"kd4"`
		CertA string `json:"certA" bson:"certA"`
		CertB string `json:"certB" bson:"certB"`
		CertC string `json:"certC" bson:"certC"`
		CertD string `json:"certD" bson:"certD"`
	}
)

type (
	Stage3 struct {
		gorm.Model
		Sk1 string `json:"sk1" bson:"sk1"`
		Sk2 string `json:"sk2" bson:"sk2"`
		Sk3 string `json:"sk3" bson:"sk3"`
		Sk4 string `json:"sk4" bson:"sk4"`
	}
)

type Stage1_Hearts struct {
	M1 string `json:"m1" bson:"m1" gorm:"unique"`
	M2 string `json:"m2" bson:"m2" gorm:"unique"`
	M3 string `json:"m3" bson:"m3" gorm:"unique"`
	M4 string `json:"m4" bson:"m4" gorm:"unique"`
}