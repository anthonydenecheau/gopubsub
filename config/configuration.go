package config

type Configuration struct {
	App        AppConfiguration
	Publisher  PublisherConfiguration
	Subscriber SubscriberConfiguration
	PubSub     PubSubConfiguration
	Logger     LoggerConfiguration
}
