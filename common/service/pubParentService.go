package service

import (
	"fmt"

	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

type parentService struct {
	pr     repository.ParentRepository
	filter string
}

type ParentService interface {
	GetFilter() string
	BuildMessage(id int64, action string) ([]*model.Parent, error)
	UpsertParent(parent *model.Parent) error
}

func (a *parentService) GetFilter() string { return a.filter }

func NewParentService(r repository.ParentRepository) ParentService {
	return &parentService{
		pr:     r,
		filter: ParentDomaine,
	}
}

func (a *parentService) BuildMessage(id int64, action string) ([]*model.Parent, error) {

	message, err := a.pr.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (a *parentService) UpsertParent(parent *model.Parent) (err error) {

	// Contrôle du timestamp ...
	d, err := a.pr.Get(parent.ID)
	if err != nil {
		return err
	}
	// ... dans le cas où le chien a été préalablement enregistré
	if d != nil {
		if d.DateMaj.After(parent.DateMaj) {
			// la date enregistrée est postérieure à la date du message
			// une mise à jour a été déjà effectuée !
			fmt.Println("Message is too old")
			return nil
		}
	}

	err = a.pr.UpsertParent(parent)
	if err != nil {
		return err
	}

	return nil
}
