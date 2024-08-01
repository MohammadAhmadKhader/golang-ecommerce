package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mohammadahmadkhader/golang-ecommerce/middlewares"
	"github.com/mohammadahmadkhader/golang-ecommerce/types"
	"github.com/mohammadahmadkhader/golang-ecommerce/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.CreateProduct).Methods("POST")
	router.HandleFunc("/products", middlewares.PaginationMiddleware(h.GetProducts)).Methods("GET")
	router.HandleFunc("/products/{id}", h.GetSingleProduct).Methods("GET")
	router.HandleFunc("/products/{id}", h.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", h.DeleteProduct).Methods("DELETE")
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var createPayload types.ProductCreatePayload
	err := utils.ParseJSON(r, &createPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(createPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	createdProd, err := h.store.CreateProduct(createPayload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "success",
		"data":    createdProd,
	})
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pagination := middlewares.GetPagination(r)
	offset := middlewares.CalculateOffset(pagination)

	products, count, err := h.store.GetProducts(pagination.Limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK,
		map[string]any{
			"products": products,
			"page":     pagination.Page,
			"limit":    pagination.Limit,
			"count":    count,
		})
}

func (h *Handler) GetSingleProduct(w http.ResponseWriter, r *http.Request) {
	var pathVars = mux.Vars(r)
	id, err := strconv.Atoi(pathVars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if id < 1 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("product id must be unsigned integer"))
		return
	}

	product, err := h.store.GetProductById(id)
	if err != nil || product.ID == 0 {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"product": product})
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var updatePayload types.ProductUpdatePayload
	err := utils.ParseJSON(r, &updatePayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	isEmpty := IsProductUpdatePayloadEmpty(updatePayload)
	if isEmpty {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("at least one field is required"))
		return
	}

	if err := utils.Validate.Struct(updatePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id must be an integer"))
		return
	}

	updatedProducts, err := h.store.UpdateProduct(id, updatePayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"message": "success",
		"data":    updatedProducts,
	})
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id must be an integer"))
		return
	}

	err = h.store.DeleteProduct(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}