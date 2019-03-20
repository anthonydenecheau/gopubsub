package repository

import (
	"database/sql"
	"log"

	model "github.com/anthonydenecheau/gopubsub/common/model"
)

type oraDogRepository struct {
	Db *sql.DB
}

// NewOraDogRepository constructor
func NewOraDogRepository(Db *sql.DB) DogRepository {
	return &oraDogRepository{Db}
}

func (oradb *oraDogRepository) GetByID(id int64, action string) ([]*model.Dog, error) {

	dogList := make([]*model.Dog, 0)

	sqlStatement := `
		SELECT 	id,
			nom,
			affixe,
			sexe,
			date_Naissance,
			pays,
			tatouage,
			transpondeur,
			codeFci,
			idRace,
			idVariete,
			race,
			variete,
			couleur,
			couleur_abr,
			code_inscription,
			id_Etalon,
			id_Lice
		FROM WS_DOG where id=:id`

	if action != "D" {
		rows, err := oradb.Db.Query(sqlStatement, id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			dog := new(model.Dog)
			err := rows.Scan(&dog.Id,
				&dog.Nom,
				&dog.Affixe,
				&dog.Sexe,
				&dog.DateNaissance,
				&dog.Pays,
				&dog.Tatouage,
				&dog.Transpondeur,
				&dog.CodeFci,
				&dog.IdRace,
				&dog.IdVariete,
				&dog.Race,
				&dog.Variete,
				&dog.Couleur,
				&dog.CouleurAbr,
				&dog.InscriptionCode,
				&dog.IdEtalon,
				&dog.IdLice)
			if err != nil {
				log.Fatal(err)
			}
			dogList = append(dogList, dog)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

	} else {
		dog := new(model.Dog)
		dog.Id = id
		dogList = append(dogList, dog)
	}

	return dogList, nil

}
