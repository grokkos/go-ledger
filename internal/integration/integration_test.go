package integration

import (
	"bytes"
	"encoding/json"
	"github.com/grokkos/go-ledger/internal/handler"
	"github.com/grokkos/go-ledger/internal/model"
	"github.com/grokkos/go-ledger/internal/server"
	"github.com/grokkos/go-ledger/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLedgerIntegration(t *testing.T) {
	// Setup test server
	ledgerService := service.NewLedgerService()
	ledgerHandler := handler.NewHandler(ledgerService)
	srv := server.NewServer(ledgerHandler)

	// Create test server
	ts := httptest.NewServer(srv.Router())
	defer ts.Close()

	// Test scenario: Deposit money and check balance
	t.Run("full transaction flow", func(t *testing.T) {
		// 1. Initial balance should be 0
		balance := getBalance(t, ts.URL)
		if balance.Amount != 0 {
			t.Errorf("Expected initial balance 0, got %v", balance.Amount)
		}

		// 2. Make a deposit
		deposit := model.TransactionRequest{
			Type:   model.Deposit,
			Amount: 100.0,
		}
		txn := createTransaction(t, ts.URL, deposit)
		if txn.Amount != 100.0 {
			t.Errorf("Expected deposit amount 100.0, got %v", txn.Amount)
		}

		// 3. Verify balance after deposit
		balance = getBalance(t, ts.URL)
		if balance.Amount != 100.0 {
			t.Errorf("Expected balance 100.0, got %v", balance.Amount)
		}

		// 4. Make a withdrawal
		withdrawal := model.TransactionRequest{
			Type:   model.Withdrawal,
			Amount: 30.0,
		}
		txn = createTransaction(t, ts.URL, withdrawal)
		if txn.Amount != 30.0 {
			t.Errorf("Expected withdrawal amount 30.0, got %v", txn.Amount)
		}

		// 5. Verify final balance
		balance = getBalance(t, ts.URL)
		if balance.Amount != 70.0 {
			t.Errorf("Expected final balance 70.0, got %v", balance.Amount)
		}

		// 6. Verify transaction history
		transactions := getTransactions(t, ts.URL)
		if len(transactions) != 2 {
			t.Errorf("Expected 2 transactions, got %d", len(transactions))
		}

		// 7. Try withdrawal with insufficient funds
		invalidWithdrawal := model.TransactionRequest{
			Type:   model.Withdrawal,
			Amount: 100.0,
		}
		code := attemptInvalidTransaction(t, ts.URL, invalidWithdrawal)
		if code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", code)
		}
	})
}

// Helper functions
func getBalance(t *testing.T, baseURL string) model.Balance {
	resp, err := http.Get(baseURL + "/api/v1/balance")
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}
	defer resp.Body.Close()

	var balance model.Balance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		t.Fatalf("Failed to decode balance: %v", err)
	}
	return balance
}

func createTransaction(t *testing.T, baseURL string, req model.TransactionRequest) model.Transaction {
	body, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/v1/transactions", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	defer resp.Body.Close()

	var txn model.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&txn); err != nil {
		t.Fatalf("Failed to decode transaction: %v", err)
	}
	return txn
}

func getTransactions(t *testing.T, baseURL string) []model.Transaction {
	resp, err := http.Get(baseURL + "/api/v1/transactions")
	if err != nil {
		t.Fatalf("Failed to get transactions: %v", err)
	}
	defer resp.Body.Close()

	var transactions []model.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		t.Fatalf("Failed to decode transactions: %v", err)
	}
	return transactions
}

func attemptInvalidTransaction(t *testing.T, baseURL string, req model.TransactionRequest) int {
	body, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/v1/transactions", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to attempt invalid transaction: %v", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
