package service

import (
	"fmt"

	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

type pedigreeService struct {
	pr     repository.PedigreeRepository
	filter string
}

type PedigreeService interface {
	GetFilter() string
	BuildMessage(id int64, action string) ([]*model.Pedigree, error)
	UpsertPedigree(pedigree *model.Pedigree) error
}

func (a *pedigreeService) GetFilter() string { return a.filter }

func NewPedigreeService(r repository.PedigreeRepository) PedigreeService {
	return &pedigreeService{
		pr:     r,
		filter: PedigreeDomaine,
	}
}

func (a *pedigreeService) BuildMessage(id int64, action string) ([]*model.Pedigree, error) {

	message, err := a.pr.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (a *pedigreeService) UpsertPedigree(pedigree *model.Pedigree) (err error) {

	// Contrôle du timestamp ...
	d, err := a.pr.Get(pedigree.ID)
	if err != nil {
		return err
	}
	// ... dans le cas où le chien a été préalablement enregistré
	if d != nil {
		if d.DateMaj.After(pedigree.DateMaj) {
			// la date enregistrée est postérieure à la date du message
			// une mise à jour a été déjà effectuée !
			fmt.Println("Message is too old")
			return nil
		}
	}

	err = a.pr.UpsertPedigree(pedigree)
	if err != nil {
		return err
	}

	return nil
}
