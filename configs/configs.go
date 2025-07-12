package configs

import (
	"Jevan/commons/appdb"
	"Jevan/commons/apploggers"
	"context"
	"fmt"
	"net/url"
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

	user := os.Getenv(MONGO_USER)
	password := os.Getenv(MONGO_PASSWORD)
	cluster := os.Getenv(MONGO_CLUSTER)

	if user == "" || password == "" || cluster == "" {
		return fmt.Errorf("missing MongoDB configuration values")
	}

	// URL-encode the password in case it contains special characters
	encodedPassword := url.QueryEscape(password)

	// Construct the URI
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=Cluster0", user, encodedPassword, cluster)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context, opts)
	if err != nil {
		logger.Errorf("MongoDB connection error: %v", err)
		return err
	}

	if err := client.Ping(context, nil); err != nil {
		logger.Errorf("MongoDB ping error: %v", err)
		return err
	}

	logger.Info("You successfully connected to MongoDB!")
	dbClient := appdb.NewDatabaseClient(os.Getenv(MONGO_DATABASE), client)
	AppConfig = &ApplicationConfig{
		HttpPort: os.Getenv(HTTP_PORT),
		DbClient: dbClient,
	}
	return nil
}
