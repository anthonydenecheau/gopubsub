package service

import (
	"fmt"

	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

type breederService struct {
	br     repository.BreederRepository
	filter string
}

type BreederService interface {
	GetAllDogs(id int64) (breeders []*model.Breeder, err error)
	GetFilter() string
	BuildMessage(id int64, action string) ([]*model.Breeder, error)
	UpsertBreeder(breeder *model.Breeder) error
}

func (a *breederService) GetFilter() string { return a.filter }

func NewBreederService(r repository.BreederRepository) BreederService {
	return &breederService{
		br:     r,
		filter: BreederDomaine,
	}
}

func (a *breederService) GetAllDogs(id int64) (breeders []*model.Breeder, err error) {
	breeders, err = a.br.GetAllDogs(id)
	if err != nil {
		return nil, err
	}

	return breeders, nil
}

func (a *breederService) BuildMessage(id int64, action string) ([]*model.Breeder, error) {

	message, err := a.br.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (a *breederService) UpsertBreeder(breeder *model.Breeder) (err error) {

	// Contrôle du timestamp ...
	d, err := a.br.Get(breeder.ID)
	if err != nil {
		return err
	}
	// ... dans le cas où le chien a été préalablement enregistré
	if d != nil {
		if d.DateMaj.After(breeder.DateMaj) {
			// la date enregistrée est postérieure à la date du message
			// une mise à jour a été déjà effectuée !
			fmt.Println("Message is too old")
			return nil
		}
	}

	err = a.br.UpsertBreeder(breeder)
	if err != nil {
		return err
	}

	return nil
}
