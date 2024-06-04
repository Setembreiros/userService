package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"userservice/cmd/provider"
	"userservice/infrastructure/atlas"
	"userservice/infrastructure/kafka"
	"userservice/internal/bus"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type app struct {
	env              string
	ctx              context.Context
	cancel           context.CancelFunc
	configuringTasks sync.WaitGroup
	runningTasks     sync.WaitGroup
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	connStr := strings.TrimSpace(os.Getenv("CONN_STR"))

	app := &app{
		env:    env,
		ctx:    ctx,
		cancel: cancel,
	}

	app.configuringLog()

	log.Info().Msgf("Starting User Service in [%s] enviroment...\n", env)

	provider := provider.NewProvider(env, connStr)

	migrator, err := provider.ProvideAtlasCLient()
	if err != nil {
		os.Exit(1)
	}
	database, err := provider.ProvideDb()
	if err != nil {
		os.Exit(1)
	}
	defer database.Client.Close()
	eventBus := provider.ProvideEventBus()
	subscriptions := provider.ProvideSubscriptions(database)
	kafkaConsumer, err := provider.ProvideKafkaConsumer(eventBus)
	if err != nil {
		os.Exit(1)
	}

	app.runConfigurationTasks(migrator, subscriptions, eventBus)
	app.runServerTasks(kafkaConsumer)
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

func (app *app) runConfigurationTasks(atlasCLient *atlas.AtlasClient, subscriptions *[]bus.EventSubscription, eventBus *bus.EventBus) {
	app.configuringTasks.Add(2)
	go app.applyMigrations(atlasCLient)
	go app.subcribeEvents(subscriptions, eventBus) // Always subscribe event before init Kafka
	app.configuringTasks.Wait()
}

func (app *app) runServerTasks(kafkaConsumer *kafka.KafkaConsumer) {
	app.runningTasks.Add(1)
	go app.initKafkaConsumption(kafkaConsumer)

	blockForever()

	app.shutdown()
}

func (app *app) applyMigrations(atlasCLient *atlas.AtlasClient) {
	defer app.configuringTasks.Done()

	err := atlasCLient.ApplyMigrations(app.ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msgf("Failed to apply migrations")
	}
}

func (app *app) subcribeEvents(subscriptions *[]bus.EventSubscription, eventBus *bus.EventBus) {
	defer app.configuringTasks.Done()

	log.Info().Msg("Subscribing events...")

	for _, subscription := range *subscriptions {
		eventBus.Subscribe(&subscription, app.ctx)
		log.Info().Msgf("%s subscribed\n", subscription.EventType)
	}

	log.Info().Msg("All events subscribed")
}

func (app *app) initKafkaConsumption(kafkaConsumer *kafka.KafkaConsumer) {
	defer app.runningTasks.Done()

	err := kafkaConsumer.InitConsumption(app.ctx)
	if err != nil {
		log.Panic().Stack().Err(err).Msg("Kafka Consumption failed")
	}
	log.Info().Msg("Kafka Consumer Group stopped")
}

func blockForever() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

func (app *app) shutdown() {
	app.cancel()
	log.Info().Msg("Shutting down Readmodels Service...")
	app.runningTasks.Wait()
	log.Info().Msg("Readmodels Service stopped")
}
