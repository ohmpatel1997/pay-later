package command

import (
	"fmt"
	"os"
	"pay-later/integration/email"
	"pay-later/integration/log"
	"pay-later/model"
	"pay-later/service/merchant"
	"pay-later/service/report"
	"pay-later/service/transaction"
	"pay-later/service/transfer"
	"pay-later/service/user"
	"strconv"
	"strings"
)

const (
	CommadCreateUser              = commadCreateUser("new user")
	CommadCreateMerchant          = commadCreateMerchant("new merchant")
	CommadCreateTransaction       = commadCreateTransaction("new txn")
	CommadUpdateMerchant          = commadUpdateMerchant("update merchant")
	CommadPayback                 = commadPayback("payback")
	CommandReportDiscount         = commandReportDiscount("report discount")
	CommandReportDues             = commandReportDues("report dues")
	CommandReportCreditLimitUsers = commandReportCreditLimitUsers("report users-at-credit-limit")
	CommandReportTotalDues        = commandReportTotalDues("report total-dues")
	CommandExit                   = commandExit("exit")
)

type CommandService interface {
	Execute(l log.Logger)
}

func NewCommand(str string) (CommandService, error) {

	if strings.HasPrefix(str, string(CommadCreateUser)) {
		return commadCreateUser(str), nil
	}

	if strings.HasPrefix(str, string(CommadCreateMerchant)) {
		return commadCreateMerchant(str), nil
	}

	if strings.HasPrefix(str, string(CommadCreateTransaction)) {
		return commadCreateTransaction(str), nil
	}

	if strings.HasPrefix(str, string(CommadUpdateMerchant)) {
		return commadUpdateMerchant(str), nil
	}

	if strings.HasPrefix(str, string(CommadPayback)) {
		return commadPayback(str), nil
	}

	if strings.HasPrefix(str, string(CommandReportDiscount)) {
		return commandReportDiscount(str), nil
	}

	if strings.HasPrefix(str, string(CommandReportDues)) {
		return commandReportDues(str), nil
	}

	if strings.HasPrefix(str, string(CommandReportCreditLimitUsers)) {
		return commandReportCreditLimitUsers(str), nil
	}

	if strings.HasPrefix(str, string(CommandReportTotalDues)) {
		return commandReportTotalDues(str), nil
	}

	if strings.HasPrefix(str, string(CommandExit)) {
		return CommandExit, nil
	}

	return nil, fmt.Errorf("invalid command")

}

type commadCreateUser string

func (c commadCreateUser) Execute(l log.Logger) {

	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	usrSrv := user.NewUserService(dbMan, emailSrv, l)

	parts := strings.Split(string(c), " ")

	name := parts[2]
	email := parts[3]
	creditLimitStr := parts[4]

	creaditLimit, err := strconv.ParseFloat(creditLimitStr, 64)
	if err != nil {
		fmt.Println("invalid limit")
		return
	}
	usr, err := usrSrv.CreateNewUser(name, email, creaditLimit)
	if err != nil {
		fmt.Println(err)
		return
	}

	creditLimitInDollars := float64(usr.CreditLimit) / float64(100)

	fmt.Println(fmt.Sprintf("%s(%0.2f)", usr.Name, creditLimitInDollars))
}

type commadCreateMerchant string

func (c commadCreateMerchant) Execute(l log.Logger) {

	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	mrtSrv := merchant.NewMerchantService(dbMan, emailSrv, l)

	parts := strings.Split(string(c), " ")

	name := parts[2]
	email := parts[3]
	discountPercentStr := strings.TrimRight(parts[4], "%")

	discountRate, err := strconv.ParseFloat(discountPercentStr, 64)
	if err != nil {
		fmt.Println("invalid limit")
		return
	}

	usr, err := mrtSrv.CreateNewMerchant(name, email, discountRate)
	if err != nil {
		fmt.Println(err)
		return
	}

	discountInDollars := float64(usr.Discount) / float64(100)

	fmt.Println(fmt.Sprintf("%s(%0.2f)", usr.Name, discountInDollars))
}

type commadCreateTransaction string

func (c commadCreateTransaction) Execute(l log.Logger) {

	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	txnSrv := transaction.NewTransactionService(dbMan, l)
	usrSrv := user.NewUserService(dbMan, emailSrv, l)
	mrtSrv := merchant.NewMerchantService(dbMan, emailSrv, l)
	transferSrv := transfer.NewTransferService(l, txnSrv, usrSrv, mrtSrv, dbMan)

	parts := strings.Split(string(c), " ")
	uname := parts[2]
	mname := parts[3]
	amount := parts[4]

	amountDollars, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println("invalid limit")
		return
	}

	_, err = transferSrv.CreateInterTransfer(uname, mname, amountDollars)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("succcess!")
}

type commadUpdateMerchant string

func (c commadUpdateMerchant) Execute(l log.Logger) {

	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	merchantSrv := merchant.NewMerchantService(dbMan, emailSrv, l)

	parts := strings.Split(string(c), " ")
	mname := parts[2]
	disStr := strings.TrimRight(parts[3], "%")

	discountRate, err := strconv.ParseFloat(disStr, 64)
	if err != nil {
		fmt.Println("invalid limit")
		return
	}

	_, err = merchantSrv.ChangeDiscountRate(mname, discountRate)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("success!")
}

type commadPayback string

func (c commadPayback) Execute(l log.Logger) {

	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	txnSrv := transaction.NewTransactionService(dbMan, l)
	usrSrv := user.NewUserService(dbMan, emailSrv, l)
	mrtSrv := merchant.NewMerchantService(dbMan, emailSrv, l)
	transferSrv := transfer.NewTransferService(l, txnSrv, usrSrv, mrtSrv, dbMan)

	parts := strings.Split(string(c), " ")
	uname := parts[1]
	amountStr := parts[2]

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("invalid amount")
		return
	}

	_, err = transferSrv.CreatePaybackTransfer(uname, amount)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("success!")
}

type commandReportDiscount string

func (c commandReportDiscount) Execute(l log.Logger) {

	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	txnSrv := transaction.NewTransactionService(dbMan, l)
	usrSrv := user.NewUserService(dbMan, emailSrv, l)
	mrtSrv := merchant.NewMerchantService(dbMan, emailSrv, l)

	rprtSrv := report.NewReportingService(l, txnSrv, usrSrv, mrtSrv, dbMan)

	parts := strings.Split(string(c), " ")
	mname := parts[2]

	dis, err := rprtSrv.GetTotalDiscount(mname)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dis)

}

type commandReportDues string

func (c commandReportDues) Execute(l log.Logger) {
	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	txnSrv := transaction.NewTransactionService(dbMan, l)
	usrSrv := user.NewUserService(dbMan, emailSrv, l)
	mrtSrv := merchant.NewMerchantService(dbMan, emailSrv, l)

	rprtSrv := report.NewReportingService(l, txnSrv, usrSrv, mrtSrv, dbMan)

	parts := strings.Split(string(c), " ")
	uname := parts[2]

	dis, err := rprtSrv.GetTotalDuesForUser(uname)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dis)
}

type commandReportCreditLimitUsers string

func (c commandReportCreditLimitUsers) Execute(l log.Logger) {
	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	txnSrv := transaction.NewTransactionService(dbMan, l)
	usrSrv := user.NewUserService(dbMan, emailSrv, l)
	mrtSrv := merchant.NewMerchantService(dbMan, emailSrv, l)

	rprtSrv := report.NewReportingService(l, txnSrv, usrSrv, mrtSrv, dbMan)

	users, err := rprtSrv.GetUsersAtCreditLimit()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(users)
}

type commandReportTotalDues string

func (c commandReportTotalDues) Execute(l log.Logger) {
	dbMan := model.NewModelManager(l)
	emailSrv := email.NewEmailService(l)
	txnSrv := transaction.NewTransactionService(dbMan, l)
	usrSrv := user.NewUserService(dbMan, emailSrv, l)
	mrtSrv := merchant.NewMerchantService(dbMan, emailSrv, l)

	rprtSrv := report.NewReportingService(l, txnSrv, usrSrv, mrtSrv, dbMan)

	str, err := rprtSrv.TotalDues()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(str)
}

type commandExit string

func (c commandExit) Execute(l log.Logger) {
	os.Exit(0)
}
