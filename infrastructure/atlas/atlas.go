package atlas

import (
	"os"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/rs/zerolog/log"
)

func NewAtlasClient() (*atlasexec.Client, error) {
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("./infrastructure/atlas/migrations"),
		),
	)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("failed to load working directory")
		return nil, err
	}
	//defer workdir.Close()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to initialize client")
		return nil, err
	}

	return client, nil
}
