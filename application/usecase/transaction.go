package usecase

import (
	"errors"

	"github.com/adelblande/codepix/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.TransactionRepositoryInterface
	PixKeyRepository model.PixKeyRepositoryInterface
}

func (t *TransactionUseCase) Register(
	accountId string, 
	amount float64, 
	pixKeyTo string, 
	pixKeyKindTo string,
	description string,
) (*model.Transaction, error) {
	account, err := t.PixKeyRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixKeyRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo)
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}

	t.TransactionRepository.Save(transaction)
	if transaction.ID == "" {
		return nil, errors.New("unable to process this transaction")
	}

	return transaction, nil
}