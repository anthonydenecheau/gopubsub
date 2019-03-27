package model

import "time"

type Owner struct {
	TableName  struct{}  `sql:"ws_dog_proprietaire" json:"-"`
	ID         int64     `json:"id"`
	Nom        string    `json:"lastName"`
	Prenom     string    `json:"firstName"`
	Adresse    string    `json:"address"`
	CodePostal string    `json:"zipCode"`
	Ville      string    `json:"town"`
	Pays       string    `json:"country"`
	IDChien    int64     `json:"idDog"`
	DateMaj    time.Time `json:"timestamp"`
}
