package models

import (
	"gorm.io/gorm"
)

type (
	Stage1 struct {
		gorm.Model
		Id     string `json:"id" bson:"id" gorm:"unique"`
		M1     string `json:"m1" bson:"m1"`
		M2     string `json:"m2" bson:"m2"`
		M3     string `json:"m3" bson:"m3"`
		M4     string `json:"m4" bson:"m4"`
		Submit bool   `json:"submit" bson:"submit"`
	}
)

type (
	Stage2 struct {
		gorm.Model
		Id     string `json:"id" bson:"id" gorm:"unique"`
		Ka1    string `json:"ka1" bson:"ka1"`
		Ka2    string `json:"ka2" bson:"ka2"`
		Ka3    string `json:"ka3" bson:"ka3"`
		Ka4    string `json:"ka4" bson:"ka4"`
		Kb1    string `json:"kb1" bson:"kb1"`
		Kb2    string `json:"kb2" bson:"kb2"`
		Kb3    string `json:"kb3" bson:"kb3"`
		Kb4    string `json:"kb4" bson:"kb4"`
		Kc1    string `json:"kc1" bson:"kc1"`
		Kc2    string `json:"kc2" bson:"kc2"`
		Kc3    string `json:"kc3" bson:"kc3"`
		Kc4    string `json:"kc4" bson:"kc4"`
		Kd1    string `json:"kd1" bson:"kd1"`
		Kd2    string `json:"kd2" bson:"kd2"`
		Kd3    string `json:"kd3" bson:"kd3"`
		Kd4    string `json:"kd4" bson:"kd4"`
		CertA  string `json:"certA" bson:"certA"`
		CertB  string `json:"certB" bson:"certB"`
		CertC  string `json:"certC" bson:"certC"`
		CertD  string `json:"certD" bson:"certD"`
		Submit bool   `json:"submit" bson:"submit"`
	}
)

type (
	Stage3 struct {
		gorm.Model
		Id  string `json:"id" bson:"id" gorm:"unique"`
		Sk1 string `json:"sk1" bson:"sk1"`
		Sk2 string `json:"sk2" bson:"sk2"`
		Sk3 string `json:"sk3" bson:"sk3"`
		Sk4 string `json:"sk4" bson:"sk4"`
	}
)

type Stage1_Hearts struct {
	M1     string `json:"m1" bson:"m1" gorm:"unique"`
	M2     string `json:"m2" bson:"m2" gorm:"unique"`
	M3     string `json:"m3" bson:"m3" gorm:"unique"`
	M4     string `json:"m4" bson:"m4" gorm:"unique"`
	S1Data string `json:"S1Data" bson:"S1Data" gorm:"unique"`
}

type Stage2_Hearts struct {
	Ka1    string `json:"Ka1" bson:"Ka1"`
	Ka2    string `json:"Ka2" bson:"Ka2"`
	Ka3    string `json:"Ka3" bson:"Ka3"`
	Ka4    string `json:"Ka4" bson:"Ka4"`
	Kb1    string `json:"Kb1" bson:"Kb1"`
	Kb2    string `json:"Kb2" bson:"Kb2"`
	Kb3    string `json:"Kb3" bson:"Kb3"`
	Kb4    string `json:"Kb4" bson:"Kb4"`
	Kc1    string `json:"Kc1" bson:"Kc1"`
	Kc2    string `json:"Kc2" bson:"Kc2"`
	Kc3    string `json:"Kc3" bson:"Kc3"`
	Kc4    string `json:"Kc4" bson:"Kc4"`
	Kd1    string `json:"Kd1" bson:"Kd1"`
	Kd2    string `json:"Kd2" bson:"Kd2"`
	Kd3    string `json:"Kd3" bson:"Kd3"`
	Kd4    string `json:"Kd4" bson:"Kd4"`
	CertA  string `json:"CertA" bson:"CertA"`
	CertB  string `json:"CertB" bson:"CertB"`
	CertC  string `json:"CertC" bson:"CertC"`
	CertD  string `json:"CertD" bson:"CertD"`
	S2Data string `json:"S2Data" bson:"S2Data"`
}

type Share_Cert struct {
	CertAs string `json:"CertAs" bson:"CertAs" gorm:"unique"`
	CertBs string `json:"CertBs" bson:"CertBs" gorm:"unique"`
	CertCs string `json:"CertCs" bson:"CertCs" gorm:"unique"`
	CertDs string `json:"CertDs" bson:"CertDs" gorm:"unique"`
}

type Fetch_Stage1 struct {
	Id string `json:"id" bson:"id" gorm:"unique"`
	M1 string `json:"m1" bson:"m1" gorm:"unique"`
	M2 string `json:"m2" bson:"m2" gorm:"unique"`
	M3 string `json:"m3" bson:"m3" gorm:"unique"`
	M4 string `json:"m4" bson:"m4" gorm:"unique"`
}

type Fetch_Stage2 struct {
	Id    string `json:"id" bson:"id" gorm:"unique"`
	Ka1   string `json:"Ka1" bson:"Ka1"`
	Ka2   string `json:"Ka2" bson:"Ka2"`
	Ka3   string `json:"Ka3" bson:"Ka3"`
	Ka4   string `json:"Ka4" bson:"Ka4"`
	Kb1   string `json:"Kb1" bson:"Kb1"`
	Kb2   string `json:"Kb2" bson:"Kb2"`
	Kb3   string `json:"Kb3" bson:"Kb3"`
	Kb4   string `json:"Kb4" bson:"Kb4"`
	Kc1   string `json:"Kc1" bson:"Kc1"`
	Kc2   string `json:"Kc2" bson:"Kc2"`
	Kc3   string `json:"Kc3" bson:"Kc3"`
	Kc4   string `json:"Kc4" bson:"Kc4"`
	Kd1   string `json:"Kd1" bson:"Kd1"`
	Kd2   string `json:"Kd2" bson:"Kd2"`
	Kd3   string `json:"Kd3" bson:"Kd3"`
	Kd4   string `json:"Kd4" bson:"Kd4"`
	CertA string `json:"CertA" bson:"CertA"`
	CertB string `json:"CertB" bson:"CertB"`
	CertC string `json:"CertC" bson:"CertC"`
	CertD string `json:"CertD" bson:"CertD"`
}

type Fetch_Cert struct {
	Id  string `json:"id" bson:"id"`
	Sk1 string `json:"sk1" bson:"sk1"`
	Sk2 string `json:"sk2" bson:"sk2"`
	Sk3 string `json:"sk3" bson:"sk3"`
	Sk4 string `json:"sk4" bson:"sk4"`
}
