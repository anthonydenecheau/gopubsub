package repository

import (
	"database/sql"
	"log"

	model "github.com/anthonydenecheau/gopubsub/common/model"
)

// BreederRepository Interfaces
type BreederRepository interface {
	GetByID(id int64, action string) ([]*model.Breeder, error)
	UpsertBreeder(breeder *model.Breeder) error
	GetAllDogs(id int64) ([]*model.Breeder, error)
	Get(id int64) (*model.Breeder, error)
}

type oraBreederRepository struct {
	Db *sql.DB
}

// NewOraBreederRepository constructor
func NewOraBreederRepository(Db *sql.DB) BreederRepository {
	return &oraBreederRepository{Db}
}

func (oradb *oraBreederRepository) Get(id int64) (*model.Breeder, error) {
	return nil, nil
}
func (oradb *oraBreederRepository) GetAllDogs(id int64) ([]*model.Breeder, error) {

	breederList := make([]*model.Breeder, 0)

	sqlStatement := `
		SELECT 	id,
			civilite,
			nom,
			prenom,
			typ_profil,
			professionnel_actif,
			raison_sociale,
			on_suffixe,
			pays,
			id_chien,
			date_maj
		FROM WS_DOG_ELEVEUR where id=:id`

	rows, err := oradb.Db.Query(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		breeder := new(model.Breeder)
		err := rows.Scan(&breeder.ID,
			&breeder.Civilite,
			&breeder.Nom,
			&breeder.Prenom,
			&breeder.TypeProfil,
			&breeder.ProfessionnelActif,
			&breeder.RaisonSociale,
			&breeder.OnSuffixe,
			&breeder.Pays,
			&breeder.IDChien,
			&breeder.DateMaj)
		if err != nil {
			log.Fatal(err)
		}
		breederList = append(breederList, breeder)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return breederList, nil
}
func (oradb *oraBreederRepository) UpsertBreeder(breeder *model.Breeder) error {
	return nil
}
func (oradb *oraBreederRepository) GetByID(id int64, action string) ([]*model.Breeder, error) {

	breederList := make([]*model.Breeder, 0)

	sqlStatement := `
		SELECT 	id,
			civilite,
			nom,
			prenom,
			typ_profil,
			professionnel_actif,
			raison_sociale,
			on_suffixe,
			pays,
			id_chien,
			date_maj
		FROM WS_DOG_ELEVEUR where id_chien=:id`

	if action != "D" {
		rows, err := oradb.Db.Query(sqlStatement, id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			breeder := new(model.Breeder)
			err := rows.Scan(&breeder.ID,
				&breeder.Civilite,
				&breeder.Nom,
				&breeder.Prenom,
				&breeder.TypeProfil,
				&breeder.ProfessionnelActif,
				&breeder.RaisonSociale,
				&breeder.OnSuffixe,
				&breeder.Pays,
				&breeder.IDChien,
				&breeder.DateMaj)
			if err != nil {
				log.Fatal(err)
			}
			breederList = append(breederList, breeder)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

	} else {
		breeder := new(model.Breeder)
		breeder.ID = id
		breederList = append(breederList, breeder)
	}

	return breederList, nil

}
