package handler

import (
	"encoding/json"
	"github.com/grokkos/go-ledger/internal/model"
	"github.com/grokkos/go-ledger/internal/service"
	"github.com/grokkos/go-ledger/pkg/response"
	"net/http"
)

type Handler struct {
	service service.LedgerService
}

func NewHandler(service service.LedgerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RecordTransaction(w http.ResponseWriter, r *http.Request) {
	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.RecordTransaction(req)
	if err != nil {
		switch err {
		case service.ErrInsufficientFunds:
			response.Error(w, err.Error(), http.StatusBadRequest)
		case service.ErrInvalidAmount:
			response.Error(w, err.Error(), http.StatusBadRequest)
		default:
			response.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response.JSON(w, transaction, http.StatusCreated)
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := h.service.GetTransactions()
	response.JSON(w, transactions, http.StatusOK)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance := h.service.GetBalance()
	response.JSON(w, balance, http.StatusOK)
}
