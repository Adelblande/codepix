package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Base `valid:"required"`
	OwnerName string `gorm:"column:owner_name;type:varchar(20);not null;" valid:"notnull"`
	Bank *Bank `valid:"-"`
	BankID string `gorm:"column:bank_id;type:uuid;not null" valid:"-"`
	Number string `gorm:"column:number;type:varchar(20);not null" valid:"notnull"`
	PixKeys []*PixKey `gorm:"ForeignKey:AccountId" valid:"-"`
}

func (account *Account) isValid() error {
	_, err := govalidator.ValidateStruct(account)
	if err != nil {
		return err
	}
	return nil
}

func NewAccount(ownerName string, bank *Bank, number string) (*Account, error){
	account := Account {
		OwnerName: ownerName,
		Bank: bank,
		Number: number,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()
	
	err := account.isValid()
	if err != nil {
		return nil, err
	}

	return &account, nil
}