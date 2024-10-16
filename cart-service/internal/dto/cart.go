package dto

import "cart-service/pkg/response"

type CartRequest struct {
	ProductId string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}

type CartDto struct {
	ProductId    string  `json:"product_id"`
	ProductImage string  `json:"product_image"`
	Name         string  `json:"name"`
	Cost         float64 `json:"cost"`
	Price        float64 `json:"price"`
	Sale         float64 `json:"sale"`
	Quantity     int64   `json:"quantity"`
	Total        float64 `json:"total"`
}

func ToCartDto(product *response.ProductResponse, quantity int64) CartDto {
	return CartDto{
		ProductId:    product.Id,
		ProductImage: product.URL,
		Cost:         product.Cost * float64(quantity),
		Name:         product.Name,
		Price:        product.Price,
		Sale:         product.Sale,
		Quantity:     quantity,
		Total:        (product.Price - (product.Price * (product.Sale / 100))) * float64(quantity),
	}
}
