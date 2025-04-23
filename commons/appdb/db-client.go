package appdb

import (
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
	d.client.Disconnect(ctx)
}

// function to get collection for the database
func (d *dbclient) Collection(collection string) DatabaseCollection {
	return newDatabaseCollection(d.client.Database(d.databaseName), collection)
}

// function to get database name
func (d *dbclient) GetDbName() string {
	return d.databaseName
}
