package merchant

import (
	"fmt"
	"pay-later/integration/email"
	"pay-later/integration/log"
	"pay-later/model"
)

type MerchantService interface {
	ChangeDiscountRate(string, float64) (*model.Merchant, error)
	GetMerchantWithName(string) (*model.Merchant, error)
	CreateNewMerchant(string, string, float64) (*model.Merchant, error)
}

type merchantService struct {
	dbSrv   model.ModelManager
	mailSrv email.EmailService
	l       log.Logger
}

func NewMerchantService(db model.ModelManager, email email.EmailService, log log.Logger) MerchantService {
	return &merchantService{
		db, email, log,
	}
}

func (u merchantService) ChangeDiscountRate(businessName string, limit float64) (*model.Merchant, error) {

	if limit < 0 || limit > 100 {
		return nil, fmt.Errorf("invalid limit")
	}

	usr := model.Merchant{
		Name: businessName,
	}

	uModel, found, err := u.dbSrv.GetWithPrimaryKey(usr)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("cannot able to find the model")
	}

	merchant, ok := uModel.(model.Merchant)
	if !ok {
		return nil, fmt.Errorf("can not able to type assert model")
	}

	merchant.Discount = getDiscountRate(limit)

	nModel, err := u.dbSrv.Upsert(merchant)
	if err != nil {
		u.l.ErrorD("can not able to update merchant", log.Fields{"user": merchant})
		return nil, err
	}

	nUser, ok := nModel.(model.Merchant)
	if !ok {
		return nil, fmt.Errorf("can not able to type assert model")
	}

	return &nUser, nil
}

func (u merchantService) CreateNewMerchant(name string, mail string, limit float64) (*model.Merchant, error) {

	if !u.mailSrv.IsValid(mail) {
		return nil, fmt.Errorf("invalid mail")
	}

	nMerchant := model.Merchant{
		Name:  name,
		Email: mail,
	}

	_, found, err := u.dbSrv.GetWithPrimaryKey(nMerchant)
	if err != nil {
		return nil, err
	}

	if found {
		return nil, fmt.Errorf("merchant already exist")
	}

	nMerchant.Discount = getDiscountRate(limit)

	uModel, err := u.dbSrv.Upsert(nMerchant)
	if err != nil {
		u.l.ErrorD("error inserting merchant into database", log.Fields{"merchant id": nMerchant.Name})
		return nil, err
	}

	merchant, ok := uModel.(model.Merchant)
	if !ok {
		return nil, fmt.Errorf("invalide model")
	}

	return &merchant, nil
}

func (u merchantService) GetMerchantWithName(name string) (*model.Merchant, error) {

	merchant := model.Merchant{
		Name: name,
	}

	nModel, found, err := u.dbSrv.GetWithPrimaryKey(merchant)
	if err != nil {
		u.l.ErrorD("can not able to get merchant with primary key", log.Fields{"primary Key": merchant.Name})
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("not found")
	}

	nUser, ok := nModel.(model.Merchant)
	if !ok {
		return nil, fmt.Errorf("can not able to type asssert user")
	}

	return &nUser, nil
}

func getDiscountRate(rate float64) int {
	return int(rate * 100)
}
