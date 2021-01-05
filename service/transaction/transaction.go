package transaction

import (
	"fmt"
	"pay-later/integration/log"
	"pay-later/model"
)

type TransactionService interface {
	CreateTransaction(*model.Transaction) (*model.Transaction, error)
	GetTotalDiscountForMerchant(string) (*int, error)
}

type transactionService struct {
	l  log.Logger
	db model.ModelManager
}

func NewTransactionService(db model.ModelManager, log log.Logger) TransactionService {
	return &transactionService{
		log, db,
	}
}

func (t transactionService) CreateTransaction(transaction *model.Transaction) (*model.Transaction, error) {

	if transaction == nil {
		t.l.ErrorD("transasction can not be nil", log.Fields{"transaction": transaction})
		return nil, fmt.Errorf("nil transaction object")
	}

	_, found, err := t.db.GetWithPrimaryKey(*transaction)
	if err != nil {
		t.l.ErrorD("error checking if the transaction already exist", log.Fields{"transaction": transaction})
		return nil, err
	}

	if found {
		return nil, fmt.Errorf("transaction already exist")
	}

	txn, err := t.db.Upsert(*transaction)
	if err != nil {
		return nil, err
	}

	nTxn := txn.(model.Transaction)

	return &nTxn, nil
}

func (t transactionService) GetTotalDiscountForMerchant(merchantName string) (*int, error) {

	var totalDiscount int = 0

	transactions, err := t.db.GetAll(model.Transaction{})
	if err != nil {
		return nil, err
	}

	for _, txn := range transactions {
		nTxn, ok := txn.(model.Transaction)
		if !ok {
			return nil, fmt.Errorf("can not able to type assert transaction")
		}

		if nTxn.Type == model.MERCHANT_DISCOUNT_CREDIT && nTxn.SourceName == merchantName && nTxn.DestinationName == model.CLEARING_ACCOUNT_NAME {
			totalDiscount += nTxn.Amount
		}
	}

	return &totalDiscount, nil
}
