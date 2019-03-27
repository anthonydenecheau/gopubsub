package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pubsub "cloud.google.com/go/pubsub"
	client "github.com/anthonydenecheau/gopubsub/common/config"
	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/service"
	"github.com/go-pg/pg"
)

// subTask is children Task
type subTask struct {
	Task
	dogService service.DogService
}

func (d subTask) Receive() error {

	ctx := context.Background()

	// REF: https://github.com/GoogleCloudPlatform/golang-samples/blob/master/pubsub/subscriptions/main.go
	// Pull messages via subscription1.
	sub := d.Task.pubService.GetClient().Subscription("dogSubscription")
	error := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message: %q\n", string(msg.Data))
		msg.Ack()

		e := new(model.Event)
		json.Unmarshal(msg.Data, &e)
		if e.Type != "Dog" {
			return
		}

		action := e.Action

		dog := new(model.Dog)
		if len(e.Message) > 0 {
			for _, msg := range e.Message {
				dog = msg
				dog.Date_maj = time.Unix(e.Timestamp/1e3, (e.Timestamp%1e3)*int64(time.Millisecond)/int64(time.Nanosecond))

				switch {
				case action == "SAVE":
				case action == "UPDATE":
					fmt.Println(">> SAVE/UPDATE event")
					err := d.dogService.UpsertDog(dog)
					if err != nil {
						fmt.Printf(">> ERROR %s\n", err)
					}
				case action == "DELETE":
					fmt.Println(">> DELETE event")
				default:
					fmt.Println(">> UNKNOWN event")
				}

				fmt.Printf("Name dog : %s\n", dog.Nom)
			}
		}
	})
	if error != nil {
		return error
	}
	return nil
}

// NewSubTask initialize all tasks
func NewSubTask(db *pg.DB, ds service.DogService, pubService client.PubSubService) {

	fmt.Println("Inside: subTask")
	task := &Task{
		pubService: pubService,
	}

	d := subTask{*task, ds}
	d.Receive()
}
