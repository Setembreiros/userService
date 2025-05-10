package provider

import (
	"userservice/infrastructure/atlas"
	"userservice/infrastructure/kafka"
	"userservice/internal/api"
	"userservice/internal/bus"
	database "userservice/internal/db"
	newuser "userservice/internal/features/new_user"
	update_userprofile "userservice/internal/features/update_useprofile"
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

func (p *Provider) ProvideEventBus() (*bus.EventBus, error) {
	kafkaProducer, err := kafka.NewKafkaProducer(p.kafkaBrokers())
	if err != nil {
		return nil, err
	}

	return bus.NewEventBus(kafkaProducer), nil
}

func (p *Provider) ProvideApiEndpoint(database *database.Database, bus *bus.EventBus) *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers(database, bus))
}

func (p *Provider) ProvideApiControllers(database *database.Database, bus *bus.EventBus) []api.Controller {
	return []api.Controller{
		update_userprofile.NewPutUserProfileController(update_userprofile.UpdateUserProfileRepository(*database), bus),
	}
}

func (p *Provider) ProvideSubscriptions(database *database.Database) *[]bus.EventSubscription {
	return &[]bus.EventSubscription{
		{
			EventType: "UserWasRegisteredEvent",
			Handler:   newuser.NewUserWasRegisteredEventHandler(newuser.NewUserRepository(*database)),
		},
	}
}

func (p *Provider) ProvideKafkaConsumer(eventBus *bus.EventBus) (*kafka.KafkaConsumer, error) {
	brokers := p.kafkaBrokers()

	return kafka.NewKafkaConsumer(brokers, eventBus)
}

func (p *Provider) kafkaBrokers() []string {
	if p.env == "development" {
		return []string{
			"localhost:9093",
		}
	} else {
		return []string{
			"172.31.0.242:9092",
			"172.31.7.110:9092",
		}
	}
}
