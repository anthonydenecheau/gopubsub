package model

import "time"

type Title struct {
	TableName     struct{}  `sql:"ws_dog_titres" json:"-"`
	ID            int64     `json:"id"`
	IDChien       int64     `json:"idDog"`
	IDTitre       int64     `json:"idTitle"`
	Code          string    `json:"title"`
	Nom           string    `json:"name"`
	Categorie     string    `json:"type"`
	Pays          string    `json:"country"`
	DateObtention string    `json:"obtentionDate"`
	DateMaj       time.Time `json:"timestamp"`
}
