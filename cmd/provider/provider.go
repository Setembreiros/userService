package provider

import (
	"userservice/infrastructure/atlas"
	"userservice/infrastructure/kafka"
	"userservice/internal/bus"
	database "userservice/internal/db"
	newuser "userservice/internal/new_user"
)

type Provider struct {
	env     string
	connStr string
}

func NewProvider(env, connStr string) *Provider {
	return &Provider{
		env:     env,
		connStr: connStr,
	}
}

func (p *Provider) ProvideAtlasCLient() (*atlas.AtlasClient, error) {
	return atlas.NewAtlasClient(p.connStr)
}

func (p *Provider) ProvideDb() (*database.Database, error) {
	return database.NewDatabase(p.connStr)
}

func (p *Provider) ProvideEventBus() *bus.EventBus {
	eventBus := bus.NewEventBus()

	return eventBus
}

func (p *Provider) ProvideSubscriptions(database *database.Database) *[]bus.EventSubscription {
	return &[]bus.EventSubscription{
		{
			EventType: "UserWasRegisteredEvent",
			Handler:   newuser.NewUserWasRegisteredEventHandler(newuser.UserProfileRepository(*database)),
		},
	}
}

func (p *Provider) ProvideKafkaConsumer(eventBus *bus.EventBus) (*kafka.KafkaConsumer, error) {
	var brokers []string

	if p.env == "development" {
		brokers = []string{
			"localhost:9093",
		}
	} else {
		brokers = []string{
			"172.31.36.175:9092",
			"172.31.45.255:9092",
		}
	}

	return kafka.NewKafkaConsumer(brokers, eventBus)
}
