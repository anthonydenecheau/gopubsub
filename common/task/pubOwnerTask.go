package task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthonydenecheau/gopubsub/common/model"
)

func (t *Task) publishOwnerChange(message []*model.Owner, action string) (serverID string, err error) {
	event := struct {
		Type      string         `json:"type"`
		Action    string         `json:"action"`
		Message   []*model.Owner `json:"message"`
		Timestamp int64          `json:"timestamp"`
	}{
		"Owner",
		action,
		message,
		makeTimestamp(),
	}

	b, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
		t.log.Errorf("   >>  error %s", err)
		return "", err
	}

	// Envoi du message
	serverID, err = t.pubService.Publish(b)
	if err != nil {
		fmt.Println(err)
		t.log.Errorf("   >>  error %s", err)
		return "", err
	}

	return serverID, err
}

func (t *Task) sendOwnerMessage(message []*model.Owner, action string) (serverID string, err error) {

	if len(message) > 0 {
		switch {
		case action == "U":
			fmt.Println(">> UPDATE event")
			t.log.Info(">> UPDATE event")
			serverID, err := t.publishOwnerChange(message, "UPDATE")
			if err != nil {
				return "", err
			}
			return serverID, nil
		case action == "I":
			fmt.Println(">> INSERT event")
			t.log.Info(">> INSERT event")
			serverID, err := t.publishOwnerChange(message, "SAVE")
			if err != nil {
				return "", err
			}
			return serverID, nil
		case action == "D":
			fmt.Println(">> DELETE event")
			t.log.Info(">> DELETE event")
			serverID, err := t.publishOwnerChange(message, "DELETE")
			if err != nil {
				return "", err
			}
			return serverID, nil
		default:
			fmt.Println(">> UNKNOWN event")
			t.log.Error(">> UNKNOWN event")
		}

	}

	return "", nil
}

func (d *pubTask) syncOwnerChanges() {

	fmt.Println("Scanning table ws_dog_sync_data : filter applied ", d.ownerService.GetFilter())
	d.log.Infof("Scanning table ws_dog_sync_data : filter applied %s", d.ownerService.GetFilter())

	ownerList, err := d.dr.GetAllChanges(d.ownerService.GetFilter())
	if err != nil {
		fmt.Println("   >> GetAllChanges error", err)
		d.log.Errorf("   >> GetAllChanges error %s", err)
		return
	}

	list := []messages{}
	// [[Boucle]] s/ le chien
	for _, owner := range ownerList {

		idDog := owner.ID
		action := owner.Action

		fmt.Println("Action required - ID", idDog)
		d.log.Infof("Action required - ID %v", idDog)

		// 1. Maj du chien, titre etc. de la table (WS_DOC_SYNC_DATA)
		fmt.Println(">> UpdateTransfert")
		d.log.Info(">> UpdateTransfert")

		err := d.dr.UpdateTransfert(idDog)
		if err != nil {
			fmt.Println("   >>  error", err)
			d.log.Errorf("   >>  error %s", err)
			continue
		}

		// 2. Lecture des infos pour le chien à synchroniser
		// Si UPDATE/INSERT et dog == null alors le chien n'est pas dans le périmètre -> on le supprime de la liste
		// + DELETE, dog == null -> on publie uniquement l'id à supprimer
		fmt.Println(">> BuildMessage")
		d.log.Info(">> BuildMessage")
		message, err := d.ownerService.BuildMessage(idDog, action)
		if err != nil {
			fmt.Println("   >>  error", err)
			d.log.Errorf("   >>  error %s", err)
			continue
		}

		if message == nil || len(message) == 0 {
			fmt.Println(">> DeleteId")
			d.log.Info(">> DeleteId")

			d.dr.DeleteId(idDog)
			if err != nil {
				fmt.Println("   >>  error", err)
				d.log.Errorf("   >>  error %s", err)
			}
			continue
		}

		// 3. Envoi du message pour maj Postgre
		/*
			fmt.Println(">> sendMessage")
			d.log.Info(">> sendMessage")

			serverID, err := d.sendMessage(message, action)
			if err != nil {
				d.log.Errorf("   >>  error %s", err)
			}
			d.log.Infof("   >>  server ID message %s", serverID)
		*/
		item := messages{owners: message, action: action}
		list = append(list, item)
	}

	// https://stackoverflow.com/questions/54612521/go-gcp-cloud-pubsub-not-batch-publishing-messages
	// publishing messages
	if len(list) > 0 {
		timeToPublish := time.Now()
		publishCount := 0
		for _, v := range list {
			fmt.Println(">> sendMessage")
			d.log.Infof(">> sendMessage %v", publishCount)

			serverID, err := d.sendOwnerMessage(v.owners, v.action)
			if err != nil {
				d.log.Errorf("   >>  error %s", err)
			}
			d.log.Infof("   >>  server ID message %s", serverID)
			publishCount++
		}
		elapsedPublish := time.Since(timeToPublish).Nanoseconds() / 1000000
		d.log.Infof("Took %v ms to publish %v messages", elapsedPublish, publishCount)
		d.log.Info("Job completed")
	}

}
