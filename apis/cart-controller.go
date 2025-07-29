package apis

import (
	"Jevan/commons"
	"Jevan/commons/apploggers"
	"Jevan/internals/models"
	"Jevan/internals/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// cartController handles operations related to cart.
type cartController struct {
	cservice services.CartService
}

// NewCartController creates a new cartController.
func NewCartController(cservice services.CartService) cartController {
	return cartController{
		cservice: cservice,
	}
}

// UpdateCart godoc
// @Summary Overwrite or add items to cart
// @Description Creates or updates a cart with new list of items and total price
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart body models.Cart true "Cart object"
// @Success 200 {object} models.Cart "Cart updated successfully"
// @Failure 400 {object} map[string]string "Invalid cart data"
// @Failure 500 {object} map[string]string "Could not update cart"
// @Router /cart [post]
func (cc *cartController) UpdateCart(e echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(e)
	var cart models.Cart
	if err := e.Bind(&cart); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cart data" + err.Error()})
	}
	cartId := e.Param("id")
	logger.Infof("Executing cart update cartId: %s, Payload: %v", cartId, cart)

	primitiveCartId, err := primitive.ObjectIDFromHex(cartId)
	if err != nil {
		logger.Error("Invalid cart ID provided")
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cart ID" + err.Error()})
	}
	cart.ID = primitiveCartId

	if err := commons.ValidateStruct(cart); err != nil {
		logger.Error("Validation failed for cart:", err)
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed for cart" + err.Error()})
	}

	logger.Infof("Executing UpdateCart %v", cart)
	if err := cc.cservice.UpdateCart(lcontext, &cart); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update cart" + err.Error()})
	}
	logger.Infof("Executed UpdateCart %v", cart)
	return e.JSON(http.StatusOK, cart)
}

// GetCartItemsById godoc
// @Summary Get all items in a cart
// @Description Get items in a cart using cartId
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} models.Cart "Cart object with all items"
// @Failure 400 {object} commons.ApiErrorResponsePayload "Failed to get items from cart"
// @Router /cart/{id} [get]
func (c *cartController) GetCartItemsById(e echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(e)
	cartId := e.Param("id")

	if len(strings.TrimSpace(cartId)) == 0 {
		logger.Error("error: cart id required.")
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse("error: cart id required.", nil))
	}

	cart, err := c.cservice.GetCartItemsById(lcontext, cartId)
	if err != nil {
		logger.Error(err)
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	logger.Info("Fetched cart items successfully")
	return e.JSON(http.StatusOK, cart)
}

// DeleteAllItems godoc
// @Summary Delete all items from cart
// @Description Remove all items from the cart identified by cartId
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200
// @Failure 400 {object} commons.ApiErrorResponsePayload "Failed to delete items from cart"
// @Router /cart/{id}/all [delete]
func (c *cartController) DeleteAllItems(e echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(e)
	cartId := e.Param("id")

	if len(strings.TrimSpace(cartId)) == 0 {
		logger.Error("error: cart id required.")
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse("error: cart id required.", nil))
	}

	if err := c.cservice.DeleteAllItems(lcontext, cartId); err != nil {
		logger.Error(err)
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	logger.Info("All items deleted from cart successfully")
	return e.NoContent(http.StatusOK)
}
