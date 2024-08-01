package cart

import (
	"fmt"

	"github.com/mohammadahmadkhader/golang-ecommerce/types"
)

// The products id's in the cart will be returned in slice.
func (h *Handler) getCartItemsIds(cart types.CartCheckoutItems) ([]int, error) {
	productsIds := make([]int, len(cart.CartItems))
	for i, cartItem := range cart.CartItems {
		if cartItem.Quantity <= 0 {
			//  i think better to return internal server error
			return nil, fmt.Errorf("product with id %v has invalid quantity", cartItem.ProductID)
		}

		productsIds[i] = cartItem.ProductID
	}

	return productsIds, nil
}

func (h *Handler) calculateTotalPrice(cartItems []types.CartCheckoutItem, productsMap map[int]types.Product) float64 {
	var totalPrice float64 = 0.0

	for _, cartItem := range cartItems {
		product := productsMap[cartItem.ProductID]
		totalPrice += product.Price
	}

	return totalPrice
}

func (h *Handler) checkProductsAvailability(cartItems []types.CartCheckoutItem, productsMap map[int]types.Product) error {
	for _, cartItem := range cartItems {
		product, ok := productsMap[cartItem.ProductID]
		if !ok {
			return fmt.Errorf("product with %v id does not exist", cartItem.ProductID)
		}

		if product.Quantity < cartItem.Quantity {
			return fmt.Errorf("you are requesting %v which is more than the available (%v)", cartItem.Quantity, product.Quantity)
		}
	}

	return nil
}

// returns productsMap[productId] => product.
func (h *Handler) createProductsMap(products []types.Product) map[int]types.Product {
	var productsMap = make(map[int]types.Product)
	for _, prod := range products {
		productsMap[prod.ID] = prod
	}
	return productsMap
}

func (h *Handler) createOrder(cartItems []types.CartCheckoutItem, productsMap map[int]types.Product, userId int) (*types.Order, error) {
	totalPrice := h.calculateTotalPrice(cartItems, productsMap)

	for _, cartItem := range cartItems {
		product := productsMap[cartItem.ProductID]
		updatePayload := types.ProductUpdatePayload{
			Quantity: product.Quantity - cartItem.Quantity,
		}

		h.productStore.UpdateProduct(product.ID, updatePayload)
	}
	
	order, err := h.orderStore.CreateOrder(types.Order{
		UserID: userId,
		Total:  totalPrice,
		Status: "pending",
		Address: "address",
	})
	if err != nil {
		return nil, err
	}

	for _, cartItem := range cartItems {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID: order.ID,
			ProductID: cartItem.ProductID,
			Quantity: cartItem.Quantity,
			Price: productsMap[cartItem.ProductID].Price,
		})
	}

	return &order, nil
}
