package model

import (
	"fmt"
	"pay-later/integration/log"
	"reflect"
)

// database tables, with indexing on `id` and `name` field
var dataBase = map[string]interface{}{
	"user":                make(map[string]User),
	"merchant":            make(map[string]Merchant),
	"transaction":         make(map[string]Transaction),
	"userpaybacktransfer": make(map[string]UserPaybackTransfer),
	"intertransfer":       make(map[string]InterTransfer),
}

type Model interface {
	TableName() string
	PrimaryKey() string
}

type ModelManager interface {
	Upsert(Model) (Model, error)
	GetWithPrimaryKey(Model) (Model, bool, error)
	GetAll(Model) ([]Model, error)
}

type modelManager struct {
	l log.Logger
}

func NewModelManager(log log.Logger) ModelManager {
	return &modelManager{
		log,
	}
}

func (m modelManager) Upsert(model Model) (Model, error) {

	s := reflect.ValueOf(model)

	if s.Kind() == reflect.Ptr {
		return nil, fmt.Errorf("please pass the non pointer to model")
	}

	primaryKey := model.PrimaryKey()
	tableName := model.TableName()

	var data interface{}
	var ok bool

	if data, ok = dataBase[tableName]; !ok {

		return nil, fmt.Errorf("invalid table name")
	}

	switch reflect.TypeOf(model) {
	case reflect.TypeOf(User{}):

		users, ok := data.(map[string]User)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}

		nUser := reflect.ValueOf(model).Convert(reflect.TypeOf(User{})).Interface().(User)

		users[primaryKey] = nUser //update the table

		dataBase[tableName] = users //update the db

		return nUser, nil

	case reflect.TypeOf(Merchant{}):

		merchants, ok := data.(map[string]Merchant)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}

		nMerchant := reflect.ValueOf(model).Convert(reflect.TypeOf(Merchant{})).Interface().(Merchant)

		merchants[primaryKey] = nMerchant

		dataBase[tableName] = merchants

		return nMerchant, nil

	case reflect.TypeOf(Transaction{}):

		transactions, ok := data.(map[string]Transaction)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}

		if _, ok := transactions[primaryKey]; ok { //transastions are append only
			return nil, fmt.Errorf("transaction already exist with given primary key")
		}

		ntransaction := reflect.ValueOf(model).Convert(reflect.TypeOf(Transaction{})).Interface().(Transaction)

		transactions[primaryKey] = ntransaction
		dataBase[tableName] = transactions

		return ntransaction, nil

	case reflect.TypeOf(InterTransfer{}):
		transfers, ok := data.(map[string]InterTransfer)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}

		if _, ok := transfers[primaryKey]; ok { //transfers are append only
			return nil, fmt.Errorf("transaction already exist with given primary key")
		}

		ntransfer := reflect.ValueOf(model).Convert(reflect.TypeOf(InterTransfer{})).Interface().(InterTransfer)

		transfers[primaryKey] = ntransfer
		dataBase[tableName] = transfers

		return ntransfer, nil

	case reflect.TypeOf(UserPaybackTransfer{}):

		transfers, ok := data.(map[string]UserPaybackTransfer)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}

		if _, ok := transfers[primaryKey]; ok { //transfers are append only
			return nil, fmt.Errorf("transaction already exist with given primary key")
		}

		ntransfer := reflect.ValueOf(model).Convert(reflect.TypeOf(UserPaybackTransfer{})).Interface().(UserPaybackTransfer)

		transfers[primaryKey] = ntransfer
		dataBase[tableName] = transfers

		return ntransfer, nil
	default:
		return nil, fmt.Errorf("invalid model type")
	}

}

func (m modelManager) GetWithPrimaryKey(model Model) (Model, bool, error) {

	s := reflect.ValueOf(model)

	if s.Kind() == reflect.Ptr {
		return nil, false, fmt.Errorf("please pass the non pointer to model")
	}

	primaryKey := model.PrimaryKey()
	tableName := model.TableName()

	var data interface{}
	var ok bool

	if data, ok = dataBase[tableName]; !ok {
		m.l.ErrorD("table name does not exist", log.Fields{"table name": model.TableName()})
		return nil, false, fmt.Errorf("invalid table name")
	}

	switch reflect.TypeOf(model) {
	case reflect.TypeOf(User{}):

		users, ok := data.(map[string]User)
		if !ok {
			return nil, false, fmt.Errorf("internal error")
		}

		if _, ok := users[primaryKey]; !ok {
			return nil, false, nil
		}

		return users[primaryKey], true, nil

	case reflect.TypeOf(Merchant{}):

		merchants, ok := data.(map[string]Merchant)
		if !ok {
			return nil, false, fmt.Errorf("internal error")
		}

		if _, ok := merchants[primaryKey]; !ok {
			return nil, false, nil
		}

		return merchants[primaryKey], true, nil

	case reflect.TypeOf(Transaction{}):

		transactions, ok := data.(map[string]Transaction)
		if !ok {
			return nil, false, fmt.Errorf("internal error")
		}

		if _, ok := transactions[primaryKey]; !ok {
			return nil, false, nil
		}

		return transactions[primaryKey], true, nil

	case reflect.TypeOf(InterTransfer{}):
		transfers, ok := data.(map[string]InterTransfer)
		if !ok {
			return nil, false, fmt.Errorf("internal error")
		}

		if _, ok := transfers[primaryKey]; !ok { //transfers are append only

			return nil, false, nil
		}

		return transfers[primaryKey], true, nil

	case reflect.TypeOf(UserPaybackTransfer{}):
		transfers, ok := data.(map[string]UserPaybackTransfer)
		if !ok {
			return nil, false, fmt.Errorf("internal error")
		}

		if _, ok := transfers[primaryKey]; !ok { //transfers are append only
			return nil, false, nil
		}

		return transfers[primaryKey], true, nil
	default:
		return nil, false, fmt.Errorf("invalid model type")
	}
}

func (m modelManager) GetAll(model Model) ([]Model, error) {
	s := reflect.ValueOf(model)

	var resp = make([]Model, 0)

	if s.Kind() == reflect.Ptr {
		return nil, fmt.Errorf("please pass the non pointer to model")
	}

	tableName := model.TableName()

	var data interface{}
	var ok bool

	if data, ok = dataBase[tableName]; !ok {
		return nil, fmt.Errorf("invalid table name")
	}

	switch reflect.TypeOf(model) {
	case reflect.TypeOf(User{}):

		users, ok := data.(map[string]User)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}

		for _, v := range users {
			resp = append(resp, v)
		}

		return resp, nil

	case reflect.TypeOf(Merchant{}):

		merchants, ok := data.(map[string]Merchant)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}

		for _, v := range merchants {
			resp = append(resp, v)
		}

		return resp, nil

	case reflect.TypeOf(Transaction{}):

		transactions, ok := data.(map[string]Transaction)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}
		for _, v := range transactions {
			resp = append(resp, v)
		}
		return resp, nil

	case reflect.TypeOf(InterTransfer{}):

		transactions, ok := data.(map[string]InterTransfer)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}
		for _, v := range transactions {
			resp = append(resp, v)
		}
		return resp, nil

	case reflect.TypeOf(UserPaybackTransfer{}):

		transactions, ok := data.(map[string]UserPaybackTransfer)
		if !ok {
			return nil, fmt.Errorf("internal error")
		}
		for _, v := range transactions {
			resp = append(resp, v)
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("invalid model type")
	}

}
