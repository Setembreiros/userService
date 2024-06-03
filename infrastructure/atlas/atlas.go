package atlas

import (
	"context"
	"errors"
	"os"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/rs/zerolog/log"
)

type AtlasClient struct {
	client  *atlasexec.Client
	workdir *atlasexec.WorkingDir
	connStr string
}

func NewAtlasClient(connStr string) (*AtlasClient, error) {
	if connStr == "" {
		log.Error().Stack().Msgf("No connection string provided")
		return nil, errors.New("no connection string provided")
	}
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("./infrastructure/atlas/migrations"),
		),
	)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("failed to load working directory")
		return nil, err
	}

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to initialize client")
		return nil, err
	}

	return &AtlasClient{
		client:  client,
		workdir: workdir,
		connStr: connStr,
	}, nil
}

func (ac *AtlasClient) ApplyMigrations(ctx context.Context) error {
	defer ac.workdir.Close()
	log.Info().Msg("Applying migrations...")

	res, err := ac.client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL: ac.connStr,
	})
	if err != nil {
		return err
	}
	log.Info().Msgf("Applied %d migrations\n", len(res.Applied))

	return nil
}
