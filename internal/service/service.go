package service

import (
	"errors"
	"fmt"
	"github.com/grokkos/go-ledger/internal/model"
	"sync"
	"time"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("amount must be positive")
)

type LedgerService interface {
	RecordTransaction(req model.TransactionRequest) (model.Transaction, error)
	GetTransactions() []model.Transaction
	GetBalance() model.Balance
}

type ledgerService struct {
	balance      float64
	transactions []model.Transaction
	mutex        sync.RWMutex
}

func (s *ledgerService) RecordTransaction(req model.TransactionRequest) (model.Transaction, error) {
	if req.Amount <= 0 {
		return model.Transaction{}, ErrInvalidAmount
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if req.Type == model.Withdrawal && s.balance < req.Amount {
		return model.Transaction{}, ErrInsufficientFunds
	}

	// Create transaction
	txn := model.Transaction{
		ID:        fmt.Sprintf("txn_%d", len(s.transactions)+1),
		Type:      req.Type,
		Amount:    req.Amount,
		Timestamp: time.Now(),
	}

	// Update balance
	if req.Type == model.Deposit {
		s.balance += req.Amount
	} else {
		s.balance -= req.Amount
	}

	s.transactions = append(s.transactions, txn)
	return txn, nil
}

func (s *ledgerService) GetTransactions() []model.Transaction {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return a copy to prevent external modifications
	transactions := make([]model.Transaction, len(s.transactions))
	copy(transactions, s.transactions)
	return transactions
}

func (s *ledgerService) GetBalance() model.Balance {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return model.Balance{Amount: s.balance}
}

func NewLedgerService() LedgerService {
	return &ledgerService{
		balance:      0,
		transactions: []model.Transaction{},
	}
}
