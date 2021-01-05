package transfer

import (
	"fmt"
	"pay-later/integration/log"
	"pay-later/model"
	"pay-later/service/merchant"
	"pay-later/service/transaction"
	"pay-later/service/user"

	"github.com/google/uuid"
)

type TransferService interface {
	CreateInterTransfer(string, string, float64) (*model.InterTransfer, error)
	CreatePaybackTransfer(string, float64) (*model.UserPaybackTransfer, error)
}

type transferService struct {
	l           log.Logger
	txnService  transaction.TransactionService
	dbSrv       model.ModelManager
	usrSrv      user.UserService
	merchantSrv merchant.MerchantService
}

func NewTransferService(l log.Logger, txnSrv transaction.TransactionService, usrSrv user.UserService, merchantSrv merchant.MerchantService, dbSrv model.ModelManager) TransferService {
	return &transferService{
		l, txnSrv, dbSrv, usrSrv, merchantSrv,
	}
}

func (t transferService) CreateInterTransfer(userName string, merchantName string, amount float64) (*model.InterTransfer, error) {

	user, err := t.usrSrv.GetUserWithName(userName)
	if err != nil || user == nil {
		return nil, err
	}

	amountToTransfer := int(amount * 100) //store it as cents

	if !user.AllowAmount(amountToTransfer) {
		return nil, fmt.Errorf("credit limit reached")
	}

	merchant, err := t.merchantSrv.GetMerchantWithName(merchantName)
	if err != nil || merchant == nil {
		return nil, err
	}

	discount := merchant.GetDiscountedAmount(amountToTransfer)
	discountedAmount := int(discount) // truncarte the 2 decimal place for cents

	actualTransferAmount := amountToTransfer - discountedAmount

	transferId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	transfer := model.InterTransfer{
		ID:             transferId,
		UserName:       user.Name,
		MerchantName:   merchant.Name,
		Amount:         amountToTransfer,
		DiscountAmount: discountedAmount,
	}

	_, found, err := t.dbSrv.GetWithPrimaryKey(transfer)
	if err != nil {
		return nil, err
	}

	if found {
		return nil, fmt.Errorf("transfer already exist")
	}

	tnsfrM, err := t.dbSrv.Upsert(transfer)
	if err != nil {
		t.l.Error("error initiating transfer", log.Fields{"transfer": transfer})
		return nil, err
	}

	nTransfer, ok := tnsfrM.(model.InterTransfer)
	if !ok {
		return nil, fmt.Errorf("can not able to type assert")
	}

	//update the user dues
	resultantDues := user.Dues + amountToTransfer
	_, err = t.usrSrv.UpdateUserDues(user.Name, resultantDues)
	if err != nil {
		return nil, err
	}

	//create transactions
	txn1Id, err := uuid.NewUUID()
	if err != nil {
		t.l.Error("error generating transaction id")
		return nil, err
	}

	txn2Id, err := uuid.NewUUID()
	if err != nil {
		t.l.Error("error generating transaction id")
		return nil, err
	}

	txn1 := model.Transaction{
		ID:              txn1Id,
		TransferID:      nTransfer.ID,
		Type:            model.USER_MERCHANT_TRANSFER,
		SourceName:      nTransfer.UserName,
		DestinationName: nTransfer.MerchantName,
		Amount:          actualTransferAmount,
	}

	txn2 := model.Transaction{
		ID:              txn2Id,
		TransferID:      nTransfer.ID,
		Type:            model.MERCHANT_DISCOUNT_CREDIT,
		SourceName:      nTransfer.MerchantName,
		DestinationName: model.CLEARING_ACCOUNT_NAME,
		Amount:          discountedAmount,
	}

	_, err = t.txnService.CreateTransaction(&txn1)
	if err != nil {
		return nil, err
	}

	_, err = t.txnService.CreateTransaction(&txn2)
	if err != nil {
		return nil, err
	}

	return &nTransfer, nil

}

func (t transferService) CreatePaybackTransfer(userName string, amount float64) (*model.UserPaybackTransfer, error) {

	user, err := t.usrSrv.GetUserWithName(userName)
	if err != nil || user == nil {
		return nil, err
	}

	amountToTransfer := int(amount * 100) //store it as cents

	if user.Dues <= 0 {
		return nil, fmt.Errorf("no dues for user")
	}

	if amountToTransfer > user.Dues {
		return nil, fmt.Errorf("payback amount should be less than or equal to dues")
	}

	finalDues := user.Dues - amountToTransfer

	nUser, err := t.usrSrv.UpdateUserDues(user.Name, finalDues)
	if err != nil {
		return nil, fmt.Errorf("can not able to change dues for user: %s", user.Name)
	}

	transferID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("can not able to generate transfer id")
	}

	transfer := model.UserPaybackTransfer{
		ID:       transferID,
		UserName: nUser.Name,
		Amount:   amountToTransfer,
	}

	_, found, err := t.dbSrv.GetWithPrimaryKey(transfer)
	if err != nil {
		return nil, fmt.Errorf("can not able to check if transfer already exist")
	}

	if found {
		return nil, fmt.Errorf("transfer already exist")
	}

	nTransfer, err := t.dbSrv.Upsert(transfer)
	if err != nil {
		return nil, err
	}

	nTrans, ok := nTransfer.(model.UserPaybackTransfer)
	if !ok {
		return nil, fmt.Errorf("can not able to type assert model")
	}

	txnID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("can not able to generate transaction id")
	}

	txn := model.Transaction{
		ID:              txnID,
		TransferID:      nTrans.ID,
		Type:            model.USER_PAYBACK_TRANSFER,
		SourceName:      model.USER_PAYBACK_ACCOUNT_NAME,
		DestinationName: nTrans.UserName,
		Amount:          amountToTransfer,
	}

	_, err = t.txnService.CreateTransaction(&txn)
	if err != nil {
		return nil, err
	}

	return &nTrans, nil
}
