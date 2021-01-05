package model

import (
	"github.com/google/uuid"
)

type TransactionType string

const (
	USER_MERCHANT_TRANSFER    = TransactionType("user-merchant")
	USER_PAYBACK_TRANSFER     = TransactionType("user-payback")
	MERCHANT_DISCOUNT_CREDIT  = TransactionType("merchant-discount")
	CLEARING_ACCOUNT_NAME     = "clearing-account"
	USER_PAYBACK_ACCOUNT_NAME = "external-account"
)

type Transaction struct {
	ID              uuid.UUID
	TransferID      uuid.UUID
	Type            TransactionType
	SourceName      string
	DestinationName string
	Amount          int //can be stored as cents
}

func (m Transaction) TableName() string {
	return "transaction"
}

func (m Transaction) PrimaryKey() string {
	return m.ID.String()
}
