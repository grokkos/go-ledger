package service

import (
	"github.com/grokkos/go-ledger/internal/model"
	"testing"
)

func TestLedgerService_RecordTransaction(t *testing.T) {
	tests := []struct {
		name         string
		request      model.TransactionRequest
		wantErr      error
		wantBalance  float64
		initialSetup func(*ledgerService)
	}{
		{
			name: "successful deposit",
			request: model.TransactionRequest{
				Type:   model.Deposit,
				Amount: 100.0,
			},
			wantErr:     nil,
			wantBalance: 100.0,
		},
		{
			name: "successful withdrawal",
			request: model.TransactionRequest{
				Type:   model.Withdrawal,
				Amount: 50.0,
			},
			wantBalance: 50.0,
			initialSetup: func(s *ledgerService) {
				s.balance = 100.0
			},
		},
		{
			name: "insufficient funds",
			request: model.TransactionRequest{
				Type:   model.Withdrawal,
				Amount: 150.0,
			},
			wantErr: ErrInsufficientFunds,
		},
		{
			name: "invalid amount",
			request: model.TransactionRequest{
				Type:   model.Deposit,
				Amount: -50.0,
			},
			wantErr: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLedgerService()
			if tt.initialSetup != nil {
				tt.initialSetup(s.(*ledgerService))
			}

			_, err := s.RecordTransaction(tt.request)
			if err != tt.wantErr {
				t.Errorf("RecordTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				balance := s.GetBalance()
				if balance.Amount != tt.wantBalance {
					t.Errorf("Balance = %v, want %v", balance.Amount, tt.wantBalance)
				}
			}
		})
	}
}
