package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/Sirupsen/logrus"
)

type pubSubService struct {
	Client       *pubsub.Client
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
	log          *logrus.Logger
}

type PubSubService interface {
	GetClient() *pubsub.Client
	GetTopicName() string
	GetTopic() *pubsub.Topic
	Publish([]byte) (string, error)
	//PublishThatScales([]byte, int) error
}

func (a *pubSubService) GetTopicName() string    { return "Topic : " + a.Topic.String() }
func (a *pubSubService) GetTopic() *pubsub.Topic { return a.Topic }

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

func (a *pubSubService) Publish(b []byte) (serverID string, err error) {

	fmt.Printf("Published a message; CountThreshold: %v\n", a.Topic.PublishSettings.CountThreshold)
	fmt.Printf("Published a message; ByteThreshold: %v\n", a.Topic.PublishSettings.ByteThreshold)
	fmt.Printf("Published a message; DelayThreshold: %v\n", a.Topic.PublishSettings.DelayThreshold)

	fmt.Println(a.GetTopic())
	fmt.Println(string(b))

	ctx := context.Background()
	res := a.Topic.Publish(ctx, &pubsub.Message{Data: b})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	serverID, err = res.Get(ctx)
	if err != nil {
		return "", err
	}
	fmt.Printf("Published a message; msg ID: %v\n", serverID)

	return serverID, nil
}

/*
func (a *pubSubService) Receive() error {

	ctx := context.Background()
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

	// Traitement des messages par lot pour équilibrer la latence et le débit
	//topic.PublishSettings.CountThreshold = 2
	topic.PublishSettings = pubsub.PublishSettings{
		ByteThreshold:  1e6, // Publish a batch when its size in bytes reaches this value. (1e6 = 1Mo)
		CountThreshold: 100, // Publish a batch when it has this many messages.
		//DelayThreshold: 10 * time.Second, // Publish a non-empty batch after this delay has passed.
		DelayThreshold: 100 * time.Millisecond,
	}

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

/*
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}
*/
