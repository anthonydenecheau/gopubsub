package repository

import (
	"fmt"

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

func (pgdb *pgDogRepository) GetByID(id int64, action string) ([]*model.Dog, error) {
	return nil, nil
}
func (pgdb *pgDogRepository) UpsertDog(dog *model.Dog) error {

	res, err := pgdb.Db.Model(dog).
		OnConflict("(Id) DO UPDATE"). // OnConflict is optional
		Insert()

	if err != nil {
		return err
	}

	fmt.Println(res, dog)

	return nil
}

func (pgdb *pgDogRepository) Get(id int64) (*model.Dog, error) {

	dog := model.Dog{ID: id}
	err := pgdb.Db.Select(&dog)

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &dog, nil
}
