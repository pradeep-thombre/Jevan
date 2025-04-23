package apis

import (
	"Jevan/commons"
	"Jevan/commons/apploggers"
	"Jevan/internals/models"
	"Jevan/internals/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type OrderController struct {
	oservice services.OrderService
}

func NewOrderController(oservice services.OrderService) *OrderController {
	return &OrderController{
		oservice: oservice,
	}
}

// @Tags Order Management
// @Summary CreateOrder
// @Description Create a new order with given details
// @Accept json
// @Produce json
// @Param payload body models.Order true "Order Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /orders [Post]
func (oc *OrderController) CreateOrder(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	logger.Info("Executing CreateOrder")
	var order *models.Order
	err := c.Bind(&order)
	if err != nil || order == nil {
		logger.Error("invalid request payload")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("invalid request payload", nil))
	}

	orderID, err := oc.oservice.CreateOrder(lcontext, order)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	logger.Infof("Executed CreateOrder, orderId: %s", orderID)
	return c.JSON(http.StatusCreated, map[string]string{
		"id": orderID,
	})
}

// @Tags Order Management
// @Summary GetOrderById
// @Description Get details of an order by its ID
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /orders/{id} [Get]
func (oc *OrderController) GetOrderById(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	orderId := c.Param("id")
	logger.Infof("Executing GetOrderById, orderId: %s", orderId)

	if len(strings.TrimSpace(orderId)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}

	order, err := oc.oservice.GetOrderById(lcontext, orderId)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	logger.Infof("Executed GetOrderById, orderId: %s", orderId)
	return c.JSON(http.StatusOK, order)
}

// @Tags Order Management
// @Summary UpdateOrder
// @Description Update an order's status or cancel the order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param payload body models.Order true "Order"
// @Success 200
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /orders/{id} [Put]
func (oc *OrderController) UpdateOrder(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	orderId := c.Param("id")
	logger.Infof("Executing UpdateOrder, orderId: %s", orderId)

	var order *models.Order
	err := c.Bind(&order)
	if err != nil || order == nil {
		logger.Error("invalid request payload")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("invalid request payload", nil))
	}

	err = oc.oservice.UpdateOrder(lcontext, orderId, order)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	logger.Infof("Executed UpdateOrder, orderId: %s", orderId)
	return c.NoContent(http.StatusOK)
}
