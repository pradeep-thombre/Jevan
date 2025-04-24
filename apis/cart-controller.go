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
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cart data"})
	}
	logger.Infof("Executing UpdateCart %v", cart)
	if err := cc.cservice.UpdateCart(lcontext, &cart); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update cart"})
	}
	logger.Infof("Executed UpdateCart %v", cart)
	return e.JSON(http.StatusOK, cart)
}

// UpdateItemQuantity godoc
// @Summary Update quantity of a cart item
// @Description Updates the quantity of a specific item in a cart, removes if quantity = 0
// @Tags Cart
// @Accept json
// @Produce json
// @Param cartId path string true "Cart ID"
// @Param itemId path string true "Item ID"
// @Param request body map[string]int true "Quantity payload" example({"quantity": 2})
// @Success 200 {object} models.Cart "Updated cart"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Failed to update item"
// @Router /cart/{cartId}/item/{itemId} [put]
func (cc *cartController) UpdateItemQuantity(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	cartId := c.Param("cartId")
	itemId := c.Param("itemId")

	type Request struct {
		Quantity int `json:"quantity"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	logger.Infof("Executing UpdateItemQuantity cart id: %s, item id", cartId, itemId)

	updatedCart, err := cc.cservice.UpdateItemQuantity(lcontext, cartId, itemId, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedCart)
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
