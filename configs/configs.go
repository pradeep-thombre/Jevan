package configs

import (
	"Jevan/commons/appdb"
	"Jevan/commons/apploggers"
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	AppConfig *ApplicationConfig
)

type ApplicationConfig struct {
	HttpPort string
	DbClient appdb.DatabaseClient
}

func NewApplicationConfig(context context.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	// Create new mongo client and connect to the server
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv(MONGO_URI)).SetServerAPIOptions(serverAPI)
	client, cerror := mongo.Connect(context, opts)
	if cerror != nil {
		logger.Errorf("Error while connecting db, error: ", cerror)
		panic(cerror)
	}

	logger.Info("You successfully connected to MongoDB!")
	dbClient := appdb.NewDatabaseClient(os.Getenv(MONGO_DATABASE), client)
	AppConfig = &ApplicationConfig{
		HttpPort: os.Getenv(HTTP_PORT),
		DbClient: dbClient,
	}
	return nil
}
