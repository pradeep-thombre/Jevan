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
// @Router /orders [post]
func (oc *OrderController) CreateOrder(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	logger.Info("Executing CreateOrder")

	var order *models.Order
	if err := c.Bind(&order); err != nil || order == nil {
		logger.Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request payload, error: "+err.Error(), nil))
	}

	if err := commons.ValidateStruct(order); err != nil {
		logger.Error("Validation failed for order:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed for order" + err.Error()})
	}

	orderID, err := oc.oservice.CreateOrder(lcontext, order)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	logger.Infof("Order Placed Successfully, orderId: %s", orderID)
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
// @Router /orders/{id} [get]
func (oc *OrderController) GetOrderById(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	orderId := c.Param("id")

	if len(strings.TrimSpace(orderId)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}

	logger.Infof("Executing GetOrderById, orderId: %s", orderId)

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
// @Router /orders/{id} [put]
func (oc *OrderController) UpdateOrder(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	orderId := c.Param("id")

	if len(strings.TrimSpace(orderId)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}

	var order *models.Order
	if err := c.Bind(&order); err != nil || order == nil {
		logger.Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request payload", nil))
	}

	logger.Infof("Executing UpdateOrder, orderId: %s", orderId)

	if err := oc.oservice.UpdateOrder(lcontext, orderId, order); err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	logger.Infof("Executed UpdateOrder, orderId: %s", orderId)
	return c.NoContent(http.StatusOK)
}

// @Tags Order Management
// @Summary GetAllOrders
// @Description Get all orders
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /orders [get]
func (oc *OrderController) GetAllOrders(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	logger.Info("Executing GetAllOrders")

	orders, err := oc.oservice.GetAllOrders(lcontext)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	response := map[string]interface{}{
		"total":  len(orders),
		"orders": orders,
	}

	logger.Infof("Executed GetAllOrders, total: %d", len(orders))
	return c.JSON(http.StatusOK, response)
}
