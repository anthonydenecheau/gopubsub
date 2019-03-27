package service

import (
	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

type personService struct {
	br     repository.BreederRepository
	or     repository.OwnerRepository
	filter string
}

type PersonService interface {
	GetFilter() string
	BuildMessageBreeder(id int64, action string) ([]*model.Breeder, error)
	//UpsertBreeder(breeder *model.Breeder) error
	BuildMessageOwner(id int64, action string) ([]*model.Owner, error)
	//UpsertOwner(owner *model.Owner) error
}

func (a *personService) GetFilter() string { return a.filter }

func NewPersonService(r repository.BreederRepository, u repository.OwnerRepository) PersonService {
	return &personService{
		br:     r,
		or:     u,
		filter: PersonDomaine,
	}
}

func (a *personService) BuildMessageBreeder(id int64, action string) ([]*model.Breeder, error) {

	message, err := a.br.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}
func (a *personService) BuildMessageOwner(id int64, action string) ([]*model.Owner, error) {

	message, err := a.or.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}
