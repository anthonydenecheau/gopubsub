package service

import (
	"fmt"

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
	UpsertDog(dog *model.Dog) error
}

func (a *dogService) GetFilter() string { return a.filter }

func NewDogService(r repository.DogRepository) DogService {
	return &dogService{
		dr:     r,
		filter: DogDomaine,
	}
}

func (a *dogService) BuildMessage(id int64, action string) ([]*model.Dog, error) {

	message, err := a.dr.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (a *dogService) UpsertDog(dog *model.Dog) (err error) {

	// Contrôle du timestamp ...
	d, err := a.dr.Get(dog.ID)
	if err != nil {
		return err
	}
	// ... dans le cas où le chien a été préalablement enregistré
	if d != nil {
		if d.Date_maj.After(dog.Date_maj) {
			// la date enregistrée est postérieure à la date du message
			// une mise à jour a été déjà effectuée !
			fmt.Println("Message is too old")
			return nil
		}
	}

	err = a.dr.UpsertDog(dog)
	if err != nil {
		return err
	}

	return nil
}
