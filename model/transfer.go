package model

import (
	"github.com/google/uuid"
)

type InterTransfer struct {
	ID             uuid.UUID
	UserName       string
	MerchantName   string
	Amount         int
	DiscountAmount int //will be store as paise, instead of rupees
}

func (m InterTransfer) TableName() string {
	return "intertransfer"
}

func (m InterTransfer) PrimaryKey() string {
	return m.ID.String()
}

type UserPaybackTransfer struct {
	ID       uuid.UUID
	UserName string
	Amount   int
}

func (m UserPaybackTransfer) TableName() string {
	return "userpaybacktransfer"
}

func (m UserPaybackTransfer) PrimaryKey() string {
	return m.ID.String()
}
