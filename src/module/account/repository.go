package account

import (
	"github.com/atletaid/go-template/src/model"
)

type AccountRepository interface {
	Create(account *model.Account) (int64, error)
	FindByID(accountID int64) (*model.Account, error)
	FindAll() (model.Accounts, error)
	Update(account *model.Account) error
}
