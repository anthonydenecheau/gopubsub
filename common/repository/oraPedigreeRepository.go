package repository

import (
	"database/sql"
	"log"

	model "github.com/anthonydenecheau/gopubsub/common/model"
)

// PedigreeRepository Interfaces
type PedigreeRepository interface {
	GetByID(id int64, action string) ([]*model.Pedigree, error)
	UpsertPedigree(pedigree *model.Pedigree) error
	Get(id int64) (*model.Pedigree, error)
}

type oraPedigreeRepository struct {
	Db *sql.DB
}

// NewOraPedigreeRepository constructor
func NewOraPedigreeRepository(Db *sql.DB) PedigreeRepository {
	return &oraPedigreeRepository{Db}
}

func (oradb *oraPedigreeRepository) Get(id int64) (*model.Pedigree, error) {
	return nil, nil
}
func (oradb *oraPedigreeRepository) UpsertPedigree(pedigree *model.Pedigree) error {
	return nil
}
func (oradb *oraPedigreeRepository) GetByID(id int64, action string) ([]*model.Pedigree, error) {

	pedigreeList := make([]*model.Pedigree, 0)

	sqlStatement := `
		SELECT 	id,
			id_chien,
			pays,
			livre,
			numero,
			date_obtention,
			date_maj
		FROM ws_dog_livres where id=:id`

	if action != "D" {
		rows, err := oradb.Db.Query(sqlStatement, id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			pedigree := new(model.Pedigree)
			err := rows.Scan(&pedigree.ID,
				&pedigree.IDChien,
				&pedigree.Pays,
				&pedigree.Livre,
				&pedigree.Numero,
				&pedigree.DateObtention,
				&pedigree.DateMaj)
			if err != nil {
				log.Fatal(err)
			}
			pedigreeList = append(pedigreeList, pedigree)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

	} else {
		pedigree := new(model.Pedigree)
		pedigree.ID = id
		pedigreeList = append(pedigreeList, pedigree)
	}

	return pedigreeList, nil

}
