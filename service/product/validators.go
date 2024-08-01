package product

import (
	"github.com/mohammadahmadkhader/golang-ecommerce/types"
)

func IsProductUpdatePayloadEmpty(payload types.ProductUpdatePayload) bool {
	if payload.Name == "" && payload.Description == "" && payload.Image == "" && payload.Price == 0 && payload.Quantity == 0 {
		return true
	}
	return false
}