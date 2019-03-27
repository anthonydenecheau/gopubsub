package repository

import (
	"database/sql"
	"log"

	model "github.com/anthonydenecheau/gopubsub/common/model"
)

// OwnerRepository Interfaces
type OwnerRepository interface {
	GetByID(id int64, action string) ([]*model.Owner, error)
	UpsertOwner(owner *model.Owner) error
	Get(id int64) (*model.Owner, error)
	GetAllDogs(id int64) ([]*model.Owner, error)
}

type oraOwnerRepository struct {
	Db *sql.DB
}

// NewOraOwnerRepository constructor
func NewOraOwnerRepository(Db *sql.DB) OwnerRepository {
	return &oraOwnerRepository{Db}
}
func (oradb *oraOwnerRepository) Get(id int64) (*model.Owner, error) {
	return nil, nil
}
func (oradb *oraOwnerRepository) GetAllDogs(id int64) ([]*model.Owner, error) {

	ownerList := make([]*model.Owner, 0)

	sqlStatement := `
		SELECT 	id,
			nom,
			prenom,
			adresse,
			code_postal,
			ville,
			pays,
			id_chien,
			date_maj
		FROM ws_dog_proprietaire where id=:id`
	rows, err := oradb.Db.Query(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		owner := new(model.Owner)
		err := rows.Scan(&owner.ID,
			&owner.Nom,
			&owner.Prenom,
			&owner.Adresse,
			&owner.CodePostal,
			&owner.Ville,
			&owner.Pays,
			&owner.IDChien,
			&owner.DateMaj)
		if err != nil {
			log.Fatal(err)
		}
		ownerList = append(ownerList, owner)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ownerList, nil
}
func (oradb *oraOwnerRepository) UpsertOwner(owner *model.Owner) error {
	return nil
}
func (oradb *oraOwnerRepository) GetByID(id int64, action string) ([]*model.Owner, error) {

	ownerList := make([]*model.Owner, 0)

	sqlStatement := `
		SELECT 	id,
			nom,
			prenom,
			adresse,
			code_postal,
			ville,
			pays,
			id_chien,
			date_maj
		FROM ws_dog_proprietaire where id_chien=:id`

	if action != "D" {
		rows, err := oradb.Db.Query(sqlStatement, id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			owner := new(model.Owner)
			err := rows.Scan(&owner.ID,
				&owner.Nom,
				&owner.Prenom,
				&owner.Adresse,
				&owner.CodePostal,
				&owner.Ville,
				&owner.Pays,
				&owner.IDChien,
				&owner.DateMaj)
			if err != nil {
				log.Fatal(err)
			}
			ownerList = append(ownerList, owner)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

	} else {
		owner := new(model.Owner)
		owner.ID = id
		ownerList = append(ownerList, owner)
	}

	return ownerList, nil

}
