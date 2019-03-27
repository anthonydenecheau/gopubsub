package model

import (
	"time"
)

type Dog struct {
	TableName        struct{}  `sql:"ws_dog" json:"-"`
	ID               int64     `json:"id"`
	Nom              string    `json:"nom"`
	Affixe           string    `json:"affixe"`
	Sexe             string    `json:"sexe"`
	DateNaissance    string    `json:"dateNaissance"`
	Pays             string    `json:"pays"`
	Tatouage         string    `json:"tatouage"`
	Transpondeur     string    `json:"transpondeur"`
	Codefci          string    `json:"codeFci"`
	Idrace           int64     `json:"idRace"`
	Idvariete        int64     `json:"idVariete"`
	Race             string    `json:"race"`
	Variete          string    `json:"variete"`
	Couleur          string    `json:"couleur"`
	Couleur_abr      string    `json:"couleurAbr"`
	Code_inscription string    `json:"inscriptionCode"`
	Id_etalon        int64     `json:"idEtalon"`
	Id_lice          int64     `json:"idLice"`
	Date_maj         time.Time `json:"timestamp"`
	On_travail       string    `json:"-"`
}
