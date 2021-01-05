package user

import (
	"fmt"
	"pay-later/integration/email"
	"pay-later/integration/log"
	"pay-later/model"
)

type UserService interface {
	ChangeCreditLimit(string, float64) (*model.User, error)
	GetUserWithName(string) (*model.User, error)
	CreateNewUser(string, string, float64) (*model.User, error)
	UpdateUserDues(string, int) (*model.User, error)
	GetCreditLimitUsers() ([]*model.User, error)
	GetTotalDues() ([]*model.User, error)
}

type userService struct {
	dbSrv   model.ModelManager
	mailSrv email.EmailService
	l       log.Logger
}

func NewUserService(db model.ModelManager, email email.EmailService, log log.Logger) UserService {
	return &userService{
		db, email, log,
	}
}

func (u userService) ChangeCreditLimit(userName string, limit float64) (*model.User, error) {

	if limit < 0 {
		return nil, fmt.Errorf("invalid limit")
	}

	usr := model.User{
		Name: userName,
	}

	uModel, found, err := u.dbSrv.GetWithPrimaryKey(usr)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("cannot able to find the model")
	}

	user, ok := uModel.(model.User)
	if !ok {
		return nil, fmt.Errorf("can not able to type assert model")
	}

	user.CreditLimit = int(limit * 100)

	nModel, err := u.dbSrv.Upsert(user)
	if err != nil {
		u.l.ErrorD("can not able to update user", log.Fields{"user": user})
		return nil, err
	}

	nUser, ok := nModel.(model.User)
	if !ok {
		return nil, fmt.Errorf("can not able to type assert model")
	}

	return &nUser, nil
}

func (u userService) CreateNewUser(name string, mail string, limit float64) (*model.User, error) {

	if !u.mailSrv.IsValid(mail) {
		return nil, fmt.Errorf("invalid mail")
	}

	nUser := model.User{
		Name:        name,
		Email:       mail,
		CreditLimit: int(limit * 100),
		Dues:        0,
	}

	_, found, err := u.dbSrv.GetWithPrimaryKey(nUser)
	if err != nil {
		return nil, err
	}

	if found {
		return nil, fmt.Errorf("user already exist")
	}

	uModel, err := u.dbSrv.Upsert(nUser)
	if err != nil {
		u.l.ErrorD("error inserting user into database", log.Fields{"user id": nUser.Name})
		return nil, err
	}

	user, ok := uModel.(model.User)
	if !ok {
		return nil, fmt.Errorf("invalide model")
	}

	return &user, nil
}

func (u userService) GetUserWithName(name string) (*model.User, error) {

	user := model.User{
		Name: name,
	}

	nModel, found, err := u.dbSrv.GetWithPrimaryKey(user)
	if err != nil {
		u.l.ErrorD("can not able to get model with primary key", log.Fields{"primary Key": user.Name})
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("not found")
	}

	nUser, ok := nModel.(model.User)
	if !ok {
		return nil, fmt.Errorf("can not able to type asssert user")
	}

	return &nUser, nil
}

func (u userService) UpdateUserDues(name string, dues int) (*model.User, error) {
	if dues < 0 {
		return nil, fmt.Errorf("dues can not be negative")
	}

	user := model.User{
		Name: name,
	}

	nModel, found, err := u.dbSrv.GetWithPrimaryKey(user)
	if err != nil {
		u.l.ErrorD("can not able to get model with primary key", log.Fields{"primary Key": user.Name})
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("user can not be found")
	}

	nUser, ok := nModel.(model.User)
	if !ok {
		return nil, fmt.Errorf("can not able to type asssert user")
	}

	if dues > nUser.CreditLimit {
		return nil, fmt.Errorf("due is over credit limit. credit limit: %d", nUser.CreditLimit)
	}

	nUser.Dues = dues

	newModel, err := u.dbSrv.Upsert(nUser)
	if err != nil {
		return nil, err
	}

	newUser, ok := newModel.(model.User)
	if !ok {
		return nil, fmt.Errorf("can not able to type assert")
	}

	return &newUser, nil
}

func (u userService) GetCreditLimitUsers() ([]*model.User, error) {
	var resp = make([]*model.User, 0)

	users, err := u.dbSrv.GetAll(model.User{})
	if err != nil {
		return resp, err
	}

	for _, mUser := range users {
		nuser, ok := mUser.(model.User)
		if !ok {
			return resp, fmt.Errorf("can not able to type assert")
		}

		if nuser.Dues >= nuser.CreditLimit {
			resp = append(resp, &nuser)
		}
	}

	return resp, nil
}

func (u userService) GetTotalDues() ([]*model.User, error) {
	var resp = make([]*model.User, 0)

	users, err := u.dbSrv.GetAll(model.User{})
	if err != nil {
		return resp, err
	}

	for _, mUser := range users {
		nuser, ok := mUser.(model.User)
		if !ok {
			return resp, fmt.Errorf("can not able to type assert")
		}

		resp = append(resp, &nuser)
	}

	return resp, nil
}
