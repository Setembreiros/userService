package provider

import (
	"context"
	"userservice/infrastructure/atlas"
	awsClients "userservice/infrastructure/aws"
	"userservice/infrastructure/kafka"
	"userservice/internal/api"
	"userservice/internal/bus"
	database "userservice/internal/db"
	newuser "userservice/internal/features/new_user"
	update_userprofile "userservice/internal/features/update_useprofile"
	objectstorage "userservice/internal/objectStorage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/rs/zerolog/log"
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

func (p *Provider) ProvideApiEndpoint(database *database.Database, objectRepository *objectstorage.ObjectStorage, bus *bus.EventBus) *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers(database, objectRepository, bus))
}

func (p *Provider) ProvideApiControllers(database *database.Database, objectRepository *objectstorage.ObjectStorage, bus *bus.EventBus) []api.Controller {
	return []api.Controller{
		update_userprofile.NewPutUserProfileController(update_userprofile.NewUpdateUserProfileRepository(database, objectRepository), bus),
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

func (p *Provider) ProvideObjectStorage(ctx context.Context) (*objectstorage.ObjectStorage, error) {
	var cfg aws.Config
	var err error

	if p.env == "development" {
		cfg, err = provideDevEnvironmentDbConfig(ctx, "4566")
	} else {
		cfg, err = provideAwsConfig(ctx)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load aws configuration")
		return nil, err
	}

	return objectstorage.NewObjectStorage(awsClients.NewS3Client(cfg, "artis-bucket")), nil
}

func (p *Provider) kafkaBrokers() []string {
	if p.env == "development" {
		return []string{
			"localhost:9093",
		}
	} else {
		return []string{
			"172.31.36.175:9092",
			"172.31.45.255:9092",
		}
	}
}

func provideAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-3"))
}

func provideDevEnvironmentDbConfig(ctx context.Context, port string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:" + port}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
}
