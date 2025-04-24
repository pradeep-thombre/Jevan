package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	"Jevan/internals/models"
	"context"
	"fmt"
)

type OrderService interface {
	CreateOrder(context context.Context, order *models.Order) (string, error)
	GetOrderById(context context.Context, orderId string) (*models.Order, error)
	UpdateOrder(context context.Context, orderId string, status *models.Order) error
}

type orderService struct {
	dbservice db.OrderDbService
}

func NewOrderService(dbservice db.OrderDbService) OrderService {
	return &orderService{
		dbservice: dbservice,
	}
}

func (os *orderService) CreateOrder(context context.Context, order *models.Order) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Info("Executing CreateOrder")

	orderID, err := os.dbservice.SaveOrder(context, order)
	if err != nil {
		logger.Error(err)
		return "", fmt.Errorf("error creating order: %s", err)
	}

	logger.Infof("Executed CreateOrder, orderId: %s", orderID)
	return orderID, nil
}

func (os *orderService) GetOrderById(context context.Context, orderId string) (*models.Order, error) {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Infof("Executing GetOrderById, orderId: %s", orderId)

	order, err := os.dbservice.GetOrderById(context, orderId)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("order not found: %s", err)
	}

	logger.Infof("Executed GetOrderById, orderId: %s", orderId)
	return order, nil
}

func (os *orderService) UpdateOrder(context context.Context, orderId string, status *models.Order) error {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Infof("Executing UpdateOrder, orderId: %s", orderId)

	err := os.dbservice.UpdateOrderStatus(context, orderId, status)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("error updating order status: %s", err)
	}

	logger.Infof("Executed UpdateOrder, orderId: %s", orderId)
	return nil
}
