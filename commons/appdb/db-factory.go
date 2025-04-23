package appdb

import (
	"context"
	"os"
	"time"

	"Jevan/commons/apploggers"
	"Jevan/commons/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppDatabaseFactory interface {
	NewDBConnection(context context.Context) (DatabaseClient, error)
	NewDbConnection(context context.Context, dbname string) (DatabaseClient, error)
}

type dbfactory struct{}

func NewDatabaseFactory() AppDatabaseFactory {
	return &dbfactory{}
}

// function to create a new connection
func (d *dbfactory) NewDBConnection(context context.Context) (DatabaseClient, error) {
	return d.NewDbConnection(context, os.Getenv(configs.MONGO_DATABASE))
}

// function to create a new connection, based on database name
func (d *dbfactory) NewDbConnection(ctx context.Context, dbname string) (DatabaseClient, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	defer logger.Sync() //nolint

	credentials := options.Credential{
		Username: os.Getenv(configs.MONGO_USER),
		Password: os.Getenv(configs.MONGO_PASSWORD),
	}

	options := options.Client().ApplyURI(os.Getenv(configs.MONGO_URI)).SetAuth(credentials)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, error := mongo.Connect(ctx, options)
	if error != nil {
		logger.Error(error)
		return nil, error
	}

	//execute ping on the connection
	if err := client.Ping(ctx, nil); err != nil {
		logger.Error(error)
		return nil, err
	}

	logger.Infof("Connected to database: " + dbname)
	return NewDatabaseClient(dbname, client), error
}
