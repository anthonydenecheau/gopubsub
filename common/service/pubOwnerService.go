package service

import (
	"fmt"

	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

type ownerService struct {
	or     repository.OwnerRepository
	filter string
}

type OwnerService interface {
	GetAllDogs(id int64) (owners []*model.Owner, err error)
	GetFilter() string
	BuildMessage(id int64, action string) ([]*model.Owner, error)
	UpsertOwner(owner *model.Owner) error
}

func (a *ownerService) GetFilter() string { return a.filter }

func NewOwnerService(r repository.OwnerRepository) OwnerService {
	return &ownerService{
		or:     r,
		filter: OwnerDomaine,
	}
}

func (a *ownerService) GetAllDogs(id int64) (owners []*model.Owner, err error) {
	owners, err = a.or.GetAllDogs(id)
	if err != nil {
		return nil, err
	}

	return owners, nil
}
func (a *ownerService) BuildMessage(id int64, action string) ([]*model.Owner, error) {

	message, err := a.or.GetByID(id, action)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (a *ownerService) UpsertOwner(owner *model.Owner) (err error) {

	// Contrôle du timestamp ...
	d, err := a.or.Get(owner.ID)
	if err != nil {
		return err
	}
	// ... dans le cas où le chien a été préalablement enregistré
	if d != nil {
		if d.DateMaj.After(owner.DateMaj) {
			// la date enregistrée est postérieure à la date du message
			// une mise à jour a été déjà effectuée !
			fmt.Println("Message is too old")
			return nil
		}
	}

	err = a.or.UpsertOwner(owner)
	if err != nil {
		return err
	}

	return nil
}
