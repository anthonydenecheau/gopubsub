package model

import "time"

type Pedigree struct {
	TableName     struct{}  `sql:"ws_dog_livres" json:"-"`
	ID            int64     `json:"id"`
	IDChien       int64     `json:"idDog"`
	Pays          string    `json:"country"`
	Livre         string    `json:"type"`
	Numero        string    `json:"number"`
	DateObtention string    `json:"obtentionDate"`
	DateMaj       time.Time `json:"timestamp"`
}
