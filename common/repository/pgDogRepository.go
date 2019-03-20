package repository

import (
	model "github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/go-pg/pg"
)

type pgDogRepository struct {
	Db *pg.DB
}

// NewPgDogRepository constructor
func NewPgDogRepository(Db *pg.DB) DogRepository {
	return &pgDogRepository{Db}
}

func (oradb *pgDogRepository) GetByID(id int64, action string) ([]*model.Dog, error) {
	return nil, nil
}
