package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type PixKeyRepositoryInterface interface {
	Registerkey(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(key string, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
	FindBank(id string) (*Bank, error)
}

type PixKey struct {
	Base `valid:"required"`
	Kind string `json:"kind" valid:"notnull"`
	Key string `json:"key" valid:"notnull"`
	AccountId string `gorm:"column:account_id;type:uuid;not null" valid:"-"`
	Status string `json:"status" valid:"notnull"`
	Account *Account `valid:"-"`
}

func (pixKey *PixKey) isValid() error {
	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid type of key")
	}
	
	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid status")
	}
	_, err := govalidator.ValidateStruct(pixKey)
	if err != nil {
		return err
	}

	return nil
}

func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey {
		Kind: kind,
		Account: account,
		Key: key,
		Status: "active",
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.isValid()
	if err != nil {
		return nil, err
	}
	
	return &pixKey, nil
}