package apis

import (
	"Jevan/commons"
	"Jevan/commons/apploggers"
	"Jevan/internals/models"
	"Jevan/internals/services"
	"net/http"

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

// AddItemToCart godoc
// @Summary Add item to cart
// @Description Add an item to the cart identified by cartId
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID" example("cart123")
// @Param item body models.CartItem true "Item to add to cart"
// @Success 200 {object} models.CartItem "Successfully added item to cart"
// @Failure 400 {object} commons.ApiErrorResponsePayload "Invalid input or failed to add item"
// @Router /cart/id [post]
func (c *cartController) AddItemToCart(e echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(e)
	cartId := e.Param("id")
	var item models.CartItem
	if err := e.Bind(&item); err != nil {
		logger.Error(err)
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	if err := c.cservice.AddItemToCart(lcontext, cartId, &item); err != nil {
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	return e.JSON(http.StatusOK, item)
}

// GetCartItemsById godoc
// @Summary Get all items in a cart
// @Description Get items in a cart using cartId
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID" example("cart123")
// @Success 200 {array} models.CartItem "List of items in cart"
// @Failure 400 {object} commons.ApiErrorResponsePayload "Failed to get items from cart"
// @Router /cart/{id} [get]
func (c *cartController) GetCartItemsById(e echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(e)
	cartId := e.Param("id")

	cart, err := c.cservice.GetCartItemsById(lcontext, cartId)
	if err != nil {
		logger.Error(err)
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	return e.JSON(http.StatusOK, cart)
}

// DeleteItemsFromCart godoc
// @Summary Delete item from cart
// @Description Remove an item from the cart using itemId
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID" example("cart123")
// @Param itemId query string true "Item ID" example("item456")
// @Success 200 {string} string "Item deleted successfully"
// @Failure 400 {object} commons.ApiErrorResponsePayload "Failed to delete item"
// @Router /cart/{id} [delete]
func (c *cartController) DeleteItemsFromCart(e echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(e)
	cartId := e.Param("id")
	itemId := e.QueryParam("itemId")

	if err := c.cservice.DeleteItemsFromCart(lcontext, cartId, itemId); err != nil {
		logger.Error(err)
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	return e.JSON(http.StatusOK, "Item deleted successfully")
}

// DeleteAllItems godoc
// @Summary Delete all items from cart
// @Description Remove all items from the cart identified by cartId
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID" example("cart123")
// @Success 200 {string} string "All items deleted successfully"
// @Failure 400 {object} commons.ApiErrorResponsePayload "Failed to delete items from cart"
// @Router /cart/{id}/all [delete]
func (c *cartController) DeleteAllItems(e echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(e)
	cartId := e.Param("id")

	if err := c.cservice.DeleteAllItems(lcontext, cartId); err != nil {
		logger.Error(err)
		return e.JSON(http.StatusBadRequest, commons.ApiErrorResponse(err.Error(), nil))
	}

	return e.JSON(http.StatusOK, "All items deleted successfully")
}
