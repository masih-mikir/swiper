package account

import (
	"log"

	"github.com/atletaid/go-template/src/model"
)

type Usecase interface {
	CreateAccount(email, fullname string) (int64, error)
	GetAccount(accountID int64) (*model.Account, error)
	GetAccounts() (model.Accounts, error)
	UpdateAccount(accountID int64, email, fullname string) error
}

type usecase struct {
	accountRepo AccountRepository
}

func NewAccountUsecase(
	accountRepo AccountRepository,
) Usecase {
	return &usecase{
		accountRepo: accountRepo,
	}
}

func (u *usecase) CreateAccount(email, fullname string) (int64, error) {
	newAccount := model.NewAccount(email, fullname)
	accountID, err := u.accountRepo.Create(newAccount)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return accountID, nil
}

func (u *usecase) GetAccount(accountID int64) (*model.Account, error) {
	account, err := u.accountRepo.FindByID(accountID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return account, nil
}

func (u *usecase) GetAccounts() (model.Accounts, error) {
	accounts, err := u.accountRepo.FindAll()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return accounts, nil
}

func (u *usecase) UpdateAccount(accountID int64, email, fullname string) error {
	account, err := u.accountRepo.FindByID(accountID)
	if err != nil {
		log.Println(err)
		return err
	}

	newAccount := *account
	newAccount.Email = email
	newAccount.Fullname = fullname

	if err = u.accountRepo.Update(&newAccount); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
