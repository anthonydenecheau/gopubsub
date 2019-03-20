package pubsub

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

type pubSubService struct {
	Client       *pubsub.Client
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}

type PubSubService interface {
	GetClient() *pubsub.Client
	GetTopic() string
	Publish([]byte) error
	//Receive() error
	//PublishThatScales([]byte, int) error
}

func (a *pubSubService) GetTopic() string { return "Topic : " + a.Topic.String() }

func (a *pubSubService) GetClient() *pubsub.Client { return a.Client }

/*
func (a *pubSubService) PublishThatScales(b []byte, n int) error {

	ctx := context.Background()

	// [START pubsub_publish_with_error_handling_that_scales]
	var wg sync.WaitGroup
	var totalErrors uint64
	t := a.Client.Topic(a.GetTopic())

	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			// The Get method blocks until a server-generated ID or
			// an error is returned for the published message.
			id, err := res.Get(ctx)
			if err != nil {
				// Error handling code can be added here.
				log.Output(1, fmt.Sprintf("Failed to publish: %v", err))
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Printf("Published message %d; msg ID: %v\n", i, id)
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return errors.New(
			fmt.Sprintf("%d of %d messages did not publish successfully",
				totalErrors, n))
	}
	return nil
	// [END pubsub_publish_with_error_handling_that_scales]
}
*/

func (a *pubSubService) Publish(b []byte) error {

	fmt.Println(a.GetTopic())
	fmt.Println(string(b))

	ctx := context.Background()
	res := a.Topic.Publish(ctx, &pubsub.Message{Data: b})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := res.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)

	return nil
}

/*
func (a *pubSubService) Receive() error {

	ctx := context.Background()

	// REF: https://github.com/GoogleCloudPlatform/golang-samples/blob/master/pubsub/subscriptions/main.go
	// Pull messages via subscription1.
	sub := a.Client.Subscription("dogSubscription")
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

//COMMENT
		var mu sync.Mutex
		received := 0
		sub := a.Client.Subscription("dogSubscription")
		cctx, cancel := context.WithCancel(ctx)
		sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			fmt.Printf("Got message: %q\n", string(msg.Data))
			mu.Lock()
			defer mu.Unlock()
			received++
			if received == 10 {
				cancel()
			}
		})
//COMMENT


	return nil
}
*/

// NewPublisher constructor
func NewPublisher() PubSubService {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, mustGetenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	topicName := mustGetenv("PUBSUB_TOPIC")
	topic := client.Topic(topicName)

	// Create the topic if it doesn't exist.
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Printf("Topic %v doesn't exist - creating it", topicName)
		_, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &pubSubService{
		Client: client,
		Topic:  topic,
	}
}

// NewSubscriber constructor
func NewSubscriber() PubSubService {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, mustGetenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	topicName := mustGetenv("PUBSUB_TOPIC")
	topic := client.Topic(topicName)

	subscriptionName := mustGetenv("PUBSUB_SUBSRCIPTION")
	subscription := client.Subscription(subscriptionName)

	// Create the topic if it doesn't exist.
	existsTop, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if !existsTop {
		log.Printf("Topic %v doesn't exist - creating it", topicName)
		_, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create the subscription if it doesn't exist.
	existsSub, err := subscription.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if !existsSub {
		log.Printf("Subscription %v doesn't exist - creating it", subscriptionName)
		_, err = client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			log.Fatal(err)
		}
	}

	return &pubSubService{
		Client:       client,
		Topic:        topic,
		Subscription: subscription,
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}
