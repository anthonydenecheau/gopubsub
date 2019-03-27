package model

import "time"

type Parent struct {
	TableName struct{}  `sql:"ws_dog_geniteurs" json:"-"`
	ID        int64     `json:"id"`
	Nom       string    `json:"name"`
	Affixe    string    `json:"affixe"`
	OnSuffixe string    `json:"onSuffixe"`
	DateMaj   time.Time `json:"timestamp"`
}
