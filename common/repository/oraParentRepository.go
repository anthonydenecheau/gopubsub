package repository

import (
	"database/sql"
	"log"

	model "github.com/anthonydenecheau/gopubsub/common/model"
)

// ParentRepository Interfaces
type ParentRepository interface {
	GetByID(id int64, action string) ([]*model.Parent, error)
	UpsertParent(parent *model.Parent) error
	Get(id int64) (*model.Parent, error)
}

type oraParentRepository struct {
	Db *sql.DB
}

// NewOraParentRepository constructor
func NewOraParentRepository(Db *sql.DB) ParentRepository {
	return &oraParentRepository{Db}
}

func (oradb *oraParentRepository) Get(id int64) (*model.Parent, error) {
	return nil, nil
}
func (oradb *oraParentRepository) UpsertParent(parent *model.Parent) error {
	return nil
}
func (oradb *oraParentRepository) GetByID(id int64, action string) ([]*model.Parent, error) {

	parentList := make([]*model.Parent, 0)

	sqlStatement := `
		SELECT 	id,
			nom,
			affixe,
			on_suffixe,
			date_maj
		FROM ws_dog_geniteurs where id=:id`

	if action != "D" {
		rows, err := oradb.Db.Query(sqlStatement, id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			parent := new(model.Parent)
			err := rows.Scan(&parent.ID,
				&parent.Nom,
				&parent.Affixe,
				&parent.OnSuffixe,
				&parent.DateMaj)
			if err != nil {
				log.Fatal(err)
			}
			parentList = append(parentList, parent)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

	} else {
		parent := new(model.Parent)
		parent.ID = id
		parentList = append(parentList, parent)
	}

	return parentList, nil

}
