package provider

import (
	"os"
	"strings"
	"userservice/infrastructure/atlas"
	"userservice/infrastructure/kafka"
	"userservice/internal/bus"
)

type Provider struct {
	env string
}

func NewProvider(env string) *Provider {
	return &Provider{
		env: env,
	}
}

func (p *Provider) ProvideAtlasCLient() (*atlas.AtlasClient, error) {
	connStr := strings.TrimSpace(os.Getenv("CONN_STR"))
	return atlas.NewAtlasClient(connStr)
}

func (p *Provider) ProvideEventBus() *bus.EventBus {
	eventBus := bus.NewEventBus()

	return eventBus
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
