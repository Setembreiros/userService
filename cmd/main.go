package main

import (
	"context"
	"os"
	"strings"
	"sync"

	"userservice/cmd/provider"
	"userservice/infrastructure/atlas"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type app struct {
	env              string
	ctx              context.Context
	cancel           context.CancelFunc
	configuringTasks sync.WaitGroup
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))

	app := &app{
		env:    env,
		ctx:    ctx,
		cancel: cancel,
	}

	app.configuringLog()

	log.Info().Msgf("Starting User Service in [%s] enviroment...\n", env)

	provider := provider.NewProvider(env)

	client, err := provider.ProvideAtlasCLient()
	if err != nil {
		os.Exit(1)
	}
	eventBus := provider.ProvideEventBus()
	_, err = provider.ProvideKafkaConsumer(eventBus)
	if err != nil {
		os.Exit(1)
	}

	app.runConfigurationTasks(client)
}

func (app *app) configuringLog() {
	if app.env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.With().Caller().Logger()
}

func (app *app) runConfigurationTasks(atlasCLient *atlas.AtlasClient) {
	app.configuringTasks.Add(1)
	go app.applyMigrations(atlasCLient)
	app.configuringTasks.Wait()
}

func (app *app) applyMigrations(atlasCLient *atlas.AtlasClient) {
	defer app.configuringTasks.Done()

	err := atlasCLient.ApplyMigrations(app.ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msgf("Failed to apply migrations")
	}
}
