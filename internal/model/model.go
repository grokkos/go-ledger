package model

import "time"

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

type Transaction struct {
	ID        string          `json:"id"`
	Type      TransactionType `json:"type"`
	Amount    float64         `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}

type Balance struct {
	Amount float64 `json:"amount"`
}

type TransactionRequest struct {
	Type   TransactionType `json:"type" validate:"required,oneof=deposit withdrawal"`
	Amount float64         `json:"amount" validate:"required,gt=0"`
}
