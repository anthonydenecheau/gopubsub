package task

import (
	"fmt"

	"github.com/anthonydenecheau/gopubsub/common/model"
)

func (d *pubTask) syncPersonChanges() {

	fmt.Println("Scanning table ws_dog_sync_data : filter applied ", d.personService.GetFilter())
	d.log.Infof("Scanning table ws_dog_sync_data : filter applied %s", d.personService.GetFilter())

	personList, err := d.dr.GetAllChanges(d.personService.GetFilter())
	if err != nil {
		fmt.Println("   >> GetAllChanges error", err)
		d.log.Errorf("   >> GetAllChanges error %s", err)
		return
	}

	// [[Boucle]] s/ la personne
	for _, person := range personList {

		idPerson := person.ID
		action := person.Action

		fmt.Println("Action required - ID", idPerson)
		d.log.Infof("Action required - ID %v", idPerson)

		// 1. Maj du chien, titre etc. de la table (WS_DOC_SYNC_DATA)
		fmt.Println(">> UpdateTransfert")
		d.log.Info(">> UpdateTransfert")

		err := d.dr.UpdateTransfert(idPerson)
		if err != nil {
			fmt.Println("   >>  error", err)
			d.log.Errorf("   >>  error %s", err)
			continue
		}

		// 2. Lecture des infos pour l'éleveur/propriétaire à synchroniser

		// PARTIE 1. Info ELEVEUR
		// Note : vue ODS_ELEVEUR (Oracle) == image de la table ODS_ELEVEUR (PostGRE)
		// Si UPDATE/INSERT et breeder == null alors l'éleveur n'est pas dans le
		// périmètre -> on le supprime de la liste
		// + DELETE, breeder == null -> on publie uniquement l'id à supprimer
		fmt.Println(">> Check Breeder")
		d.log.Info(">> Check Breeder")
		breeders := make([]*model.Breeder, 0)
		if action != "D" {
			breeders, err = d.breederService.GetAllDogs(idPerson)
			if err != nil {
				fmt.Println("   >>  error", err)
				d.log.Errorf("   >>  error %s", err)
			}
		} else {
			breeder := new(model.Breeder)
			breeder.ID = idPerson
			breeders = append(breeders, breeder)
		}

		// Envoi du message
		if len(breeders) > 0 {
			if breeders[0].ID > 0 {
				d.log.Info(">> sendMessage")

				serverID, err := d.sendBreederMessage(breeders, action)
				if err != nil {
					d.log.Errorf("   >>  error %s", err)
				}
				d.log.Infof("   >>  server ID message %s", serverID)
			} else {
				fmt.Println(">> No Breeder Found")
				d.log.Info(">> No Breeder Foundr")
			}
		} else {
			fmt.Println(">> No Breeder Found")
			d.log.Info(">> No Breeder Foundr")
		}

		// PARTIE 2. Info PROPRIETAIRE
		// Note : vue ODS_PROPRIETAIRE (Oracle) == image de la table ODS_PROPRIETAIRE
		// (PostGRE)
		// Si UPDATE/INSERT et owner == null alors le propriétaire n'est pas dans le
		// périmètre -> on le supprime de la liste
		// + DELETE, owner == null -> on publie uniquement l'id à supprimer
		fmt.Println(">> Check Owner")
		d.log.Info(">> Check Owner")
		owners := make([]*model.Owner, 0)
		if action != "D" {
			owners, err = d.ownerService.GetAllDogs(idPerson)
			if err != nil {
				fmt.Println("   >>  error", err)
				d.log.Errorf("   >>  error %s", err)
			}
		} else {
			owner := new(model.Owner)
			owner.ID = idPerson
			owners = append(owners, owner)
		}

		// Envoi du message
		if len(owners) > 0 {
			if owners[0].ID > 0 {
				fmt.Println(">> sendMessage")
				d.log.Info(">> sendMessage")

				serverID, err := d.sendOwnerMessage(owners, action)
				if err != nil {
					d.log.Errorf("   >>  error %s", err)
				}
				d.log.Infof("   >>  server ID message %s", serverID)
			} else {
				fmt.Println(">> No Owner Found")
				d.log.Info(">> No Owner Foundr")
			}
		} else {
			fmt.Println(">> No Owner Found")
			d.log.Info(">> No Owner Foundr")
		}

		if (owners == nil || len(owners) == 0) && (breeders == nil || len(breeders) == 0) {
			fmt.Println(">> DeleteId")
			d.log.Info(">> DeleteId")

			d.dr.DeleteId(idPerson)
		}
	}

}
