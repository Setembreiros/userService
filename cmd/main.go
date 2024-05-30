package main

import (
	"os"
	"strings"

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
