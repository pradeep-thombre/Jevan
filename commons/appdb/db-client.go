package appdb

import (
	"Jevan/commons/apploggers"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseClient interface {
	GetDbName() string
	Disconnect(ctx context.Context)
	Collection(collection string) DatabaseCollection
}

type dbclient struct {
	databaseName string
	client       *mongo.Client
}

func NewDatabaseClient(databasename string, client *mongo.Client) DatabaseClient {
	return &dbclient{
		databaseName: databasename,
		client:       client,
	}
}

// function to get close the db connection
func (d *dbclient) Disconnect(ctx context.Context) {
	logger := apploggers.GetLogger(ctx, false)
	if err := d.client.Disconnect(ctx); err != nil {
		logger.Info("Error disconnecting from DB: %v\n", err)
	}
}

// function to get collection for the database
func (d *dbclient) Collection(collection string) DatabaseCollection {
	return newDatabaseCollection(d.client.Database(d.databaseName), collection)
}

// function to get database name
func (d *dbclient) GetDbName() string {
	return d.databaseName
}
