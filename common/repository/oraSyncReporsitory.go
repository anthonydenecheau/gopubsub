package repository

import (
	"database/sql"
	"log"

	model "github.com/anthonydenecheau/gopubsub/common/model"
)

type oraSyncRepository struct {
	Db *sql.DB
}

// SyncRepository Interfaces
type SyncRepository interface {
	GetAllChanges(domaine string) ([]*model.SyncData, error)
	UpdateTransfert(ID int64) error
	DeleteId(ID int64) error
}

// NewOraSyncRepository constructor
func NewOraSyncRepository(Db *sql.DB) SyncRepository {
	return &oraSyncRepository{Db}
}

func (oradb *oraSyncRepository) UpdateTransfert(ID int64) error {

	sqlStatement := `
		UPDATE WS_DOG_SYNC_DATA
		SET on_transfert = 'O'
		WHERE id = :id`

	res, err := oradb.Db.Exec(sqlStatement, ID)
	_, err = res.RowsAffected()
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		return err
	default:
		return nil
	}
}

func (oradb *oraSyncRepository) DeleteId(ID int64) error {

	sqlStatement := `
		DELETE WS_DOG_SYNC_DATA
		WHERE id = :id`

	res, err := oradb.Db.Exec(sqlStatement, ID)
	_, err = res.RowsAffected()
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		return err
	default:
		return nil
	}
}

func (oradb *oraSyncRepository) GetAllChanges(domaine string) ([]*model.SyncData, error) {

	rows, err := oradb.Db.Query("select id, action, on_transfert from WS_DOG_SYNC_DATA where on_transfert='N' and domaine=:domaine", domaine)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dogList := make([]*model.SyncData, 0)
	for rows.Next() {
		syncDog := new(model.SyncData)
		err := rows.Scan(&syncDog.ID, &syncDog.Action, &syncDog.Transfert)
		if err != nil {
			log.Fatal(err)
		}
		dogList = append(dogList, syncDog)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dogList, nil
}
