package repository

import model "github.com/anthonydenecheau/gopubsub/common/model"

// DogRepository Interfaces
type DogRepository interface {
	GetByID(id int64, action string) ([]*model.Dog, error)
	UpsertDog(dog *model.Dog) error
	Get(id int64) (*model.Dog, error)
}
