package handler

import (
	"bytes"
	"encoding/json"
	"github.com/grokkos/go-ledger/internal/model"
	"github.com/grokkos/go-ledger/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockLedgerService is a mock implementation of LedgerService
type MockLedgerService struct {
	transactions []model.Transaction
	balance      float64
	recordErr    error
}

func (m *MockLedgerService) RecordTransaction(req model.TransactionRequest) (model.Transaction, error) {
	if m.recordErr != nil {
		return model.Transaction{}, m.recordErr
	}
	txn := model.Transaction{
		ID:     "test_txn",
		Type:   req.Type,
		Amount: req.Amount,
	}
	return txn, nil
}

func (m *MockLedgerService) GetTransactions() []model.Transaction {
	return m.transactions
}

func (m *MockLedgerService) GetBalance() model.Balance {
	return model.Balance{Amount: m.balance}
}

func TestHandler_RecordTransaction(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    *MockLedgerService
		wantStatus     int
		wantErrMessage string
	}{
		{
			name: "successful deposit",
			requestBody: model.TransactionRequest{
				Type:   model.Deposit,
				Amount: 100.0,
			},
			mockService: &MockLedgerService{},
			wantStatus:  http.StatusCreated,
		},
		{
			name:        "invalid request body",
			requestBody: "invalid json",
			mockService: &MockLedgerService{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: model.TransactionRequest{
				Type:   model.Withdrawal,
				Amount: 100.0,
			},
			mockService: &MockLedgerService{
				recordErr: service.ErrInsufficientFunds,
			},
			wantStatus:     http.StatusBadRequest,
			wantErrMessage: "insufficient funds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with mock service
			h := NewHandler(tt.mockService)

			// Create request
			var body bytes.Buffer
			if err := json.NewEncoder(&body).Encode(tt.requestBody); err != nil {
				t.Fatal(err)
			}

			// Create test request and response recorder
			req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", &body)
			rec := httptest.NewRecorder()

			// Handle request
			h.RecordTransaction(rec, req)

			// Check status code
			if rec.Code != tt.wantStatus {
				t.Errorf("RecordTransaction() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			// Check error message if applicable
			if tt.wantErrMessage != "" {
				var response map[string]string
				if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				if response["error"] != tt.wantErrMessage {
					t.Errorf("RecordTransaction() error = %v, want %v", response["error"], tt.wantErrMessage)
				}
			}
		})
	}
}
