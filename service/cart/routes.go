package cart

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mohammadahmadkhader/golang-ecommerce/service/auth"
	"github.com/mohammadahmadkhader/golang-ecommerce/types"
	"github.com/mohammadahmadkhader/golang-ecommerce/utils"
)

type Handler struct {
	productStore types.ProductStore
	orderStore   types.OrderStore
	userStore    types.UserStore
}

func NewHandler(productStore types.ProductStore, orderStore types.OrderStore, userStore types.UserStore) *Handler {
	return &Handler{
		productStore: productStore,
		orderStore:   orderStore,
		userStore:    userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.AuthenticationMiddleware(h.handleCheckout)).Methods("POST")
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	var cart types.CartCheckoutItems
	err := utils.ParseJSON(r, &cart)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate.Struct(cart)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	tokenPayload, err := auth.GetTokenPayload(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	productsIds, err := h.getCartItemsIds(cart)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := h.productStore.GetProductsByID(productsIds)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	productsMap := h.createProductsMap(products)
	err = h.checkProductsAvailability(cart.CartItems, productsMap)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userId := tokenPayload.UserId
	order, err := h.createOrder(cart.CartItems,productsMap,userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"order":order,
	})
}
