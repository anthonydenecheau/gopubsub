package service

import (
	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

type dogService struct {
	dr     repository.DogRepository
	filter string
}

type DogService interface {
	GetFilter() string
	BuildMessage(id int64, action string) ([]*model.Dog, error)
}

func (a *dogService) GetFilter() string { return a.filter }

func NewDogService(r repository.DogRepository) DogService {
	return &dogService{
		dr:     r,
		filter: "CHIEN",
	}
}

func (a *dogService) BuildMessage(id int64, action string) ([]*model.Dog, error) {

	message, err := a.dr.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}
