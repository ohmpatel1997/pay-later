package report

import (
	"fmt"
	"pay-later/integration/log"
	"pay-later/model"
	"pay-later/service/merchant"
	"pay-later/service/transaction"
	"pay-later/service/user"
)

type ReportService interface {
	GetTotalDiscount(string) (string, error)
	GetTotalDuesForUser(string) (string, error)
	GetUsersAtCreditLimit() ([]string, error)
	TotalDues() (string, error)
}

type reportService struct {
	l           log.Logger
	txnService  transaction.TransactionService
	dbSrv       model.ModelManager
	usrSrv      user.UserService
	merchantSrv merchant.MerchantService
}

func NewReportingService(l log.Logger, txnSrv transaction.TransactionService, usrSrv user.UserService, merchantSrv merchant.MerchantService, dbSrv model.ModelManager) ReportService {
	return &reportService{
		l, txnSrv, dbSrv, usrSrv, merchantSrv,
	}
}

func (r reportService) GetTotalDiscount(merchant string) (string, error) {
	dues, err := r.txnService.GetTotalDiscountForMerchant(merchant)
	if err != nil {
		return "", err
	}
	if dues == nil {
		return "", fmt.Errorf("can not able to find dues")
	}

	amountInDollars := float64(*dues) / float64(100)
	return fmt.Sprintf("%0.2f", amountInDollars), nil
}

func (r reportService) GetTotalDuesForUser(name string) (string, error) {
	usr, err := r.usrSrv.GetUserWithName(name)
	if err != nil {
		return "", err
	}

	duesInDollars := float64(usr.Dues) / float64(100)

	return fmt.Sprintf("%0.2f", duesInDollars), nil
}

func (r reportService) GetUsersAtCreditLimit() ([]string, error) {

	var resp = make([]string, 0)
	users, err := r.usrSrv.GetCreditLimitUsers()
	if err != nil {
		return resp, err
	}

	for _, usr := range users {
		if usr != nil {
			resp = append(resp, usr.Name)
		}
	}

	return resp, nil
}

func (r reportService) TotalDues() (string, error) {
	users, err := r.usrSrv.GetTotalDues()
	if err != nil {
		return "", err
	}

	var resp = ""
	var total float64 = 0
	for _, usr := range users {
		if usr != nil {

			duesInDollars := float64(usr.Dues) / float64(100)
			total += duesInDollars
			resp += fmt.Sprintf("%s: %0.2f\n", usr.Name, duesInDollars)
		}

	}

	resp += fmt.Sprintf("total: %0.2f", total)
	return resp, nil
}
