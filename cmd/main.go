package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"userservice/cmd/provider"
	"userservice/infrastructure/atlas"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type app struct {
	env string
}

func main() {
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))

	app := &app{
		env: env,
	}

	app.configuringLog()

	log.Info().Msgf("Starting User Service in [%s] enviroment...\n", env)

	provider := provider.NewProvider(env)

	client, err := atlas.NewAtlasClient()
	if err != nil {
		log.Fatal().Stack().Err(err).Msgf("Failed to create Atlas client")
	}

	// Run `atlas migrate apply` on a SQLite database under /tmp.
	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable",
	})
	if err != nil {
		log.Fatal().Msgf("failed to apply migrations: %v", err)
	}
	fmt.Printf("Applied %d migrations\n", len(res.Applied))

	eventBus := provider.ProvideEventBus()
	_, err = provider.ProvideKafkaConsumer(eventBus)
	if err != nil {
		os.Exit(1)
	}
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
