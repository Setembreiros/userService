package provider

import (
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

func (p *Provider) ProvideEventBus() *bus.EventBus {
	eventBus := bus.NewEventBus()

	return eventBus
}
