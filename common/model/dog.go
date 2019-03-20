package model

import (
	"time"
)

type Dog struct {
	TableName        struct{} `sql:"ws_dog"`
	ID               int64
	Nom              string
	Affixe           string
	Sexe             string
	DateNaissance    string `json:"DateNaissance"`
	Pays             string
	Tatouage         string
	Transpondeur     string
	Codefci          string
	Idrace           int64
	Idvariete        int64
	Race             string
	Variete          string
	Couleur          string
	Couleur_abr      string
	Code_inscription string
	Id_etalon        int64
	Id_lice          int64
	Date_maj         time.Time
	On_travail       string
}
