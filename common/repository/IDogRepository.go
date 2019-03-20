package repository

import model "github.com/anthonydenecheau/gopubsub/common/model"

// DogRepository Interfaces
type DogRepository interface {
	GetByID(id int64, action string) ([]*model.Dog, error)
}
