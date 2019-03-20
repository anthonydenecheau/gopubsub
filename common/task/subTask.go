package task

import (
	"context"
	"encoding/json"
	"fmt"

	pubsub "cloud.google.com/go/pubsub"
	"github.com/anthonydenecheau/gopubsub/common/model"
	client "github.com/anthonydenecheau/gopubsub/common/pubsub"
	"github.com/anthonydenecheau/gopubsub/common/service"
	"github.com/go-pg/pg"
)

// subTask is children Task
type subTask struct {
	Task
	dogService service.DogService
}

func (t *Task) Receive() error {

	ctx := context.Background()

	// REF: https://github.com/GoogleCloudPlatform/golang-samples/blob/master/pubsub/subscriptions/main.go
	// Pull messages via subscription1.
	sub := t.pubService.GetClient().Subscription("dogSubscription")
	error := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message: %q\n", string(msg.Data))
		msg.Ack()

		e := new(model.Event)
		json.Unmarshal(msg.Data, &e)
		dog := new(model.Dog)
		dog = e.Message
		fmt.Printf("Name dog : %s\n", dog.Nom)
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
