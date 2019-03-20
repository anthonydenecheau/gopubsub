package task

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/pubsub"
	"github.com/anthonydenecheau/gopubsub/common/service"

	syncRepository "github.com/anthonydenecheau/gopubsub/common/repository"
)

func (t Task) pubTopic() string { return t.pubService.GetTopic() }

// PubTask is children Task
type pubTask struct {
	Task
	dogService service.DogService
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (t Task) publishChange(message []*model.Dog, action string) {
	event := new(model.Event)
	event.Type = "Dog"
	event.Action = action
	// le sub Java n'accepte pas d'Array
	event.Message = message[0]
	event.Timestamp = makeTimestamp()

	b, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Envoi du message
	t.pubService.Publish(b)

}

func (t Task) sendMessage(message []*model.Dog, action string) {

	if len(message) > 0 {
		switch {
		case action == "U":
			fmt.Println(">> UPDATE event")
			t.publishChange(message, action)
		case action == "I":
			fmt.Println(">> INSERT event")
		case action == "D":
			fmt.Println(">> DELETE event")
		default:
			fmt.Println(">> UNKNOWN event")
		}

	}
	return
}

// a behavior only available for the PubTask
func (d pubTask) syncChanges() {

	fmt.Println("Scanning table ws_dog_sync_data : filter applied ", d.dogService.GetFilter())

	dogList, err := d.dr.GetAllChanges(d.dogService.GetFilter())
	if err != nil {
		fmt.Println(">> GetAllChanges error", err)
		return
	}

	// [[Boucle]] s/ le chien
	for _, dog := range dogList {

		idDog := dog.ID
		action := dog.Action

		fmt.Println("Action to ID", idDog)

		// 1. Maj du chien, titre etc. de la table (WS_DOC_SYNC_DATA)
		fmt.Println(">> UpdateTransfert")
		err := d.dr.UpdateTransfert(idDog)
		if err != nil {
			fmt.Println("	>>  error", err)
			continue
		}

		// 2. Lecture des infos pour le chien à synchroniser
		// Si UPDATE/INSERT et dog == null alors le chien n'est pas dans le périmètre -> on le supprime de la liste
		// + DELETE, dog == null -> on publie uniquement l'id à supprimer
		fmt.Println(">> BuildMessage")
		message, err := d.dogService.BuildMessage(idDog, action)
		if err != nil {
			fmt.Println("	>>  error", err)
			continue
		}

		if message == nil || len(message) == 0 {
			fmt.Println(">> DeleteId")
			d.dr.DeleteId(idDog)
			if err != nil {
				fmt.Println("	>>  error", err)
			}
			continue
		}

		// 3. Envoi du message pour maj Postgre
		fmt.Println(">> sendMessage")
		d.sendMessage(message, action)

	}
}

// NewTask initialize all tasks
func NewTask(db *sql.DB, ds service.DogService, pubService pubsub.PubSubService) {

	fmt.Println("Inside: pubTask")
	dr := syncRepository.NewOraSyncRepository(db)

	task := &Task{
		dr:         dr,
		pubService: pubService,
	}

	d := pubTask{*task, ds}
	fmt.Println("PubTopic: ", d.pubTopic()) //has Task  behavior
	d.syncChanges()                         //has DogTask behavior

}
