package repository

import (
	"database/sql"
	"log"

	model "github.com/anthonydenecheau/gopubsub/common/model"
)

// TitleRepository Interfaces
type TitleRepository interface {
	GetByID(id int64, action string) ([]*model.Title, error)
	UpsertTitle(title *model.Title) error
	Get(id int64) (*model.Title, error)
}

type oraTitleRepository struct {
	Db *sql.DB
}

// NewOraTitleRepository constructor
func NewOraTitleRepository(Db *sql.DB) TitleRepository {
	return &oraTitleRepository{Db}
}

func (oradb *oraTitleRepository) Get(id int64) (*model.Title, error) {
	return nil, nil
}
func (oradb *oraTitleRepository) UpsertTitle(title *model.Title) error {
	return nil
}
func (oradb *oraTitleRepository) GetByID(id int64, action string) ([]*model.Title, error) {

	titleList := make([]*model.Title, 0)

	sqlStatement := `
		SELECT 	id,
			id_chien,
			id_titre,
			code,
			nom,
			categorie,
			pays,
			date_obtention,
			date_maj
		FROM ws_dog_titres where id=:id`

	if action != "D" {
		rows, err := oradb.Db.Query(sqlStatement, id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			title := new(model.Title)
			err := rows.Scan(&title.ID,
				&title.IDChien,
				&title.IDTitre,
				&title.Code,
				&title.Nom,
				&title.Categorie,
				&title.Pays,
				&title.DateObtention,
				&title.DateMaj)
			if err != nil {
				log.Fatal(err)
			}
			titleList = append(titleList, title)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

	} else {
		title := new(model.Title)
		title.ID = id
		titleList = append(titleList, title)
	}

	return titleList, nil

}
