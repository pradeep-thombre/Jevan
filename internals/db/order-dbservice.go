package db

import (
	"Jevan/commons/appdb"
	"Jevan/commons/apploggers"
	"Jevan/internals/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DbService interface {
	SaveOrder(ctx context.Context, order *models.Order) (string, error)
	GetOrderById(ctx context.Context, orderId string) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderId string, status *models.Order) error
}

type orderDbService struct {
	ucollection appdb.DatabaseCollection
}

func NewOrderDbService(dbclient appdb.DatabaseClient) DbService {
	return &orderDbService{
		ucollection: dbclient.Collection("orders"),
	}
}

func (o *orderDbService) SaveOrder(ctx context.Context, order *models.Order) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Info("Executing SaveOrder")

	result, err := o.ucollection.InsertOne(ctx, order)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	orderId := result.InsertedID.(primitive.ObjectID).Hex()
	logger.Infof("Executed SaveOrder, orderId: %s", orderId)
	return orderId, nil
}

func (o *orderDbService) GetOrderById(ctx context.Context, orderId string) (*models.Order, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetOrderById, orderId: %s", orderId)

	id, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		return nil, fmt.Errorf("invalid orderId: %s", orderId)
	}

	var order *models.Order
	filter := bson.M{"_id": id}
	err = o.ucollection.FindOne(ctx, filter, &order)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Infof("Executed GetOrderById, orderId: %s", orderId)
	return order, nil
}

func (o *orderDbService) UpdateOrderStatus(ctx context.Context, orderId string, status *models.Order) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing UpdateOrderStatus, orderId: %s", orderId)

	id, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		return fmt.Errorf("invalid orderId: %s", orderId)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status.Status, "updated_at": status.UpdatedAt}}

	_, err = o.ucollection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("Executed UpdateOrderStatus, orderId: %s", orderId)
	return nil
}
