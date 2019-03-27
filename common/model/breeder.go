package model

import "time"

type Breeder struct {
	TableName          struct{}  `sql:"ws_dog_eleveur" json:"-"`
	ID                 int64     `json:"id"`
	Civilite           string    `json:"civility"`
	Nom                string    `json:"lastName"`
	Prenom             string    `json:"firstName"`
	TypeProfil         string    `json:"typeProfil"`
	ProfessionnelActif string    `json:"professionnelActif"`
	RaisonSociale      string    `json:"raisonSociale"`
	Pays               string    `json:"pays"`
	OnSuffixe          string    `json:"onSuffixe"`
	IDChien            int64     `json:"idDog"`
	DateMaj            time.Time `json:"timestamp"`
}
