package service

import (
	"fmt"

	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

type titleService struct {
	tr     repository.TitleRepository
	filter string
}

type TitleService interface {
	GetFilter() string
	BuildMessage(id int64, action string) ([]*model.Title, error)
	UpsertTitle(title *model.Title) error
}

func (a *titleService) GetFilter() string { return a.filter }

func NewTitleService(r repository.TitleRepository, filter string) TitleService {
	return &titleService{
		tr:     r,
		filter: filter,
	}
}

func (a *titleService) BuildMessage(id int64, action string) ([]*model.Title, error) {

	message, err := a.tr.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (a *titleService) UpsertTitle(title *model.Title) (err error) {

	// Contrôle du timestamp ...
	d, err := a.tr.Get(title.ID)
	if err != nil {
		return err
	}
	// ... dans le cas où le chien a été préalablement enregistré
	if d != nil {
		if d.DateMaj.After(title.DateMaj) {
			// la date enregistrée est postérieure à la date du message
			// une mise à jour a été déjà effectuée !
			fmt.Println("Message is too old")
			return nil
		}
	}

	err = a.tr.UpsertTitle(title)
	if err != nil {
		return err
	}

	return nil
}
