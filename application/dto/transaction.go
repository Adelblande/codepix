package dto

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
)

type TransactionDto struct {
	ID string `json:"id" valid:"required;uuid"`
	AccountId string `json:"accountId" valid:"required;uuid"`
	Amount float64 `json:"amount" valid:"required;numeric"`
	PixKeyTo string `json:"pixKeyTo" valid:"required"`
	PixKeyKindTo string `json:"pixKeyKindTo" valid:"required"`
	Description string `json:"description" valid:"required"`
	Status string `json:"status" valid:"required"`
	Error string `json:"error" valid:"required"`
}

func (t *TransactionDto) isValid() error {
	_, err := govalidator.ValidateStruct(t)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionDto) ParseJson(dataJson []byte) error {
	err := json.Unmarshal(dataJson, t)
	if err != nil {
		return err
	}

	err = t.isValid()
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionDto) ToJson() ([]byte, error) {
	err := t.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewTransactionDto() *TransactionDto {
	return &TransactionDto{}
}