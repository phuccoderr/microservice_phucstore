package controller

import (
	"cart-service/global"
	"cart-service/internal/constants"
	"cart-service/internal/dto"
	"cart-service/internal/middleware"
	"cart-service/internal/service"
	"cart-service/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CartController struct {
	cartService    service.ICartService
	productService service.IProductService
}

func NewCartController(cartService service.ICartService, productService service.IProductService) *CartController {
	return &CartController{
		cartService:    cartService,
		productService: productService,
	}
}

func (cc *CartController) AddProductToCart(c *gin.Context) {
	var cartRequest *dto.CartRequest
	err := c.ShouldBindJSON(&cartRequest)
	if err != nil {
		global.Logger.Info("Should Bind Json AddProductToCart", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_BAD_REQUEST, http.StatusBadRequest)
		return
	}

	_, customer, err := middleware.JWTGetTokenAndCustomer(c)
	if err != nil {
		global.Logger.Error("Get Token Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_UNAUTHORIZED, http.StatusUnauthorized)
		return
	}

	newQuantity, err := cc.cartService.AddProductToCart(customer.ID, cartRequest.ProductId, cartRequest.Quantity)
	if err != nil {
		global.Logger.Error("Add To Cart Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_BAD_REQUEST, http.StatusBadRequest)
		return
	}

	response.SuccessResponse(c, http.StatusCreated, constants.CREATE_SUCCESS, newQuantity)

}

func (cc *CartController) GetCart(c *gin.Context) {

	token, customer, err := middleware.JWTGetTokenAndCustomer(c)
	if err != nil {
		global.Logger.Error("Get Token Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_UNAUTHORIZED, http.StatusUnauthorized)
		return
	}

	cartsRedis, err := cc.cartService.GetCart(customer.ID)
	if err != nil {
		global.Logger.Error("Get Cart Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	var carts []dto.CartDto
	for _, cart := range cartsRedis {
		product, err := cc.productService.GetProductById(cart.ProductId, token)
		if err != nil {
			global.Logger.Error("Get Product Error", zap.Error(err))
			response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}
		carts = append(carts, dto.ToCartDto(product, customer.ID, cart.Quantity))
	}

	response.SuccessResponse(c, http.StatusOK, constants.GET_SUCCESS, carts)

}

func (cc *CartController) DeleteCart(c *gin.Context) {
	productId := c.Param("id")

	if productId == "" {
		global.Logger.Error("ProductId is empty")
		response.ErrorResponse(c, "productId is required", constants.STATUS_BAD_REQUEST, http.StatusBadRequest)
		return
	}

	_, customer, err := middleware.JWTGetTokenAndCustomer(c)
	if err != nil {
		global.Logger.Error("Get Token Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_UNAUTHORIZED, http.StatusUnauthorized)
		return
	}

	cart, err := cc.cartService.CheckProductInCart(customer.ID, productId)
	if err != nil {
		global.Logger.Error("Check Product Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	if !cart {
		global.Logger.Error("Cart not exists!")
		response.ErrorResponse(c, constants.DB_NOT_FOUND, constants.STATUS_NOT_FOUND, http.StatusNotFound)
		return
	}

	err = cc.cartService.DeleteCart(customer.ID, productId)
	if err != nil {
		global.Logger.Error("Delete Cart Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, http.StatusOK, constants.DELETE_SUCCESS, nil)
}

func (cc *CartController) Checkout(c *gin.Context) {
	token, customer, err := middleware.JWTGetTokenAndCustomer(c)
	if err != nil {
		global.Logger.Error("Get Token Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_UNAUTHORIZED, http.StatusUnauthorized)
		return
	}

	cartsRedis, err := cc.cartService.GetCart(customer.ID)
	if err != nil {
		global.Logger.Error("Get Cart Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	var carts []dto.CartDto
	for _, cart := range cartsRedis {
		product, err := cc.productService.GetProductById(cart.ProductId, token)
		if err != nil {
			global.Logger.Error("Get Product Error", zap.Error(err))
			response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}
		if cart.Quantity > product.Stock {
			global.Logger.Info(fmt.Sprintf("quantity less than %v", product.Stock))
			response.ErrorResponse(c, fmt.Sprintf("quantity less than %v", product.Stock),
				constants.STATUS_STATUS_UNPROCESSABLEENTITTY,
				http.StatusUnprocessableEntity)
			return
		}
		carts = append(carts, dto.ToCartDto(product, customer.ID, cart.Quantity))
	}
	if len(carts) == 0 {
		global.Logger.Info("Cart not exists!")
		response.ErrorResponse(c, constants.DB_NOT_FOUND, constants.STATUS_NOT_FOUND, http.StatusNotFound)
		return
	}

	checkOutInfo := cc.cartService.Checkout(carts)
	response.SuccessResponse(c, http.StatusOK, constants.CHECKOUT_SUCCESS, checkOutInfo)
}

func (cc *CartController) PlaceOrder(c *gin.Context) {
	placeOrder := &dto.PlaceOrderRequest{}
	err := c.ShouldBindJSON(&placeOrder)
	if err != nil {
		global.Logger.Info("Bind Json Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_BAD_REQUEST, http.StatusBadRequest)
		return
	}

	token, customer, err := middleware.JWTGetTokenAndCustomer(c)
	if err != nil {
		global.Logger.Error("Get Token Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_UNAUTHORIZED, http.StatusUnauthorized)
		return
	}

	cartsRedis, err := cc.cartService.GetCart(customer.ID)
	if err != nil {
		global.Logger.Error("Get Cart Error", zap.Error(err))
		response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	orderMessage := dto.PlaceOrderMessage{}

	var carts []dto.CartDto
	for _, cart := range cartsRedis {
		product, err := cc.productService.GetProductById(cart.ProductId, token)
		if err != nil {
			global.Logger.Error("Get Product Error", zap.Error(err))
			response.ErrorResponse(c, err.Error(), constants.STATUS_INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}
		if cart.Quantity > product.Stock {
			global.Logger.Info(fmt.Sprintf("quantity less than %v", product.Stock))
			response.ErrorResponse(c, fmt.Sprintf("quantity less than %v", product.Stock),
				constants.STATUS_STATUS_UNPROCESSABLEENTITTY,
				http.StatusUnprocessableEntity)
			return
		}
		carts = append(carts, dto.ToCartDto(product, customer.ID, cart.Quantity))
	}

	if len(carts) == 0 {
		global.Logger.Info("Cart not exists!")
		response.ErrorResponse(c, constants.DB_NOT_FOUND, constants.STATUS_NOT_FOUND, http.StatusNotFound)
		return
	}

	checkOutInfo := cc.cartService.Checkout(carts)

	orderMessage.CheckOut = checkOutInfo
	orderMessage.CustomerId = customer.ID
	orderMessage.CustomerEmail = customer.Email
	orderMessage.Items = carts
	orderMessage.Address = placeOrder.Address
	orderMessage.PaymentMethod = placeOrder.PaymentMethod
	orderMessage.PhoneNumber = placeOrder.PhoneNumber

	response.SuccessResponse(c, http.StatusOK, constants.PLACE_ORDER_SUCCESS, orderMessage)
}