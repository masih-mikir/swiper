package repository

import (
	"encoding/json"
	"fmt"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	cache "github.com/patrickmn/go-cache"
	"github.com/sportivaid/go-template/src/common/apperror"
	"github.com/sportivaid/go-template/src/model"
	"github.com/sportivaid/go-template/src/module/account"
)

const (
	KeyAccountsFindAll = "accounts:find_all"
	KeyAccountsFind    = "accounts:find"
)

type redisAccountRepo struct {
	cache map[string]*cache.Cache
	pool  *redigo.Pool
	next  account.AccountRepository
}

func NewMiddlewareAccountRepository(cache map[string]*cache.Cache, pool *redigo.Pool, next account.AccountRepository) account.AccountRepository {
	return &redisAccountRepo{
		cache: cache,
		pool:  pool,
		next:  next,
	}
}

func (repo *redisAccountRepo) do(command string, args ...interface{}) (reply interface{}, err error) {
	conn := repo.pool.Get()
	defer conn.Close()

	return conn.Do(command, args...)
}

func (repo *redisAccountRepo) clearAllFindListCache() error {
	keys := []interface{}{
		KeyAccountsFindAll,
	}

	if _, err := repo.do("DEL", keys...); err != nil {
		log.Println(err)
		return apperror.InternalServerError
	}

	repo.cache[KeyAccountsFindAll].Flush()
	return nil
}

func (repo *redisAccountRepo) Create(account *model.Account) (int64, error) {
	lastID, err := repo.next.Create(account)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	if err := repo.clearAllFindListCache(); err != nil {
		log.Println(err)
		return 0, err
	}

	return lastID, nil
}

func (repo *redisAccountRepo) FindByID(accountID int64) (*model.Account, error) {
	field := fmt.Sprintf("%v", accountID)

	if accountCache, found := repo.cache[KeyAccountsFind].Get(field); found {
		return accountCache.(*model.Account), nil
	}

	accountJSON, err := redigo.Bytes(repo.do("HGET", KeyAccountsFind, field))
	if err != nil {
		account, err := repo.next.FindByID(accountID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		accountJSON, err := json.Marshal(&account)
		if err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := repo.do("HSET", KeyAccountsFind, field, accountJSON); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := repo.do("EXPIRE", KeyAccountsFind, 3600); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		repo.cache[KeyAccountsFind].SetDefault(field, account)
		return account, nil
	}

	var account *model.Account
	if err := json.Unmarshal(accountJSON, &account); err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	repo.cache[KeyAccountsFind].SetDefault(field, account)
	return account, nil
}

func (repo *redisAccountRepo) FindAll() (model.Accounts, error) {
	field := "*"

	if accountCache, found := repo.cache[KeyAccountsFindAll].Get(field); found {
		return accountCache.(model.Accounts), nil
	}

	accountsJSON, err := redigo.Bytes(repo.do("HGET", KeyAccountsFindAll, field))
	if err != nil {
		accounts, err := repo.next.FindAll()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		accountsJSON, err := json.Marshal(&accounts)
		if err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := repo.do("HSET", KeyAccountsFindAll, field, accountsJSON); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		if _, err := repo.do("EXPIRE", KeyAccountsFindAll, 3600); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		repo.cache[KeyAccountsFindAll].SetDefault(field, accounts)
		return accounts, nil
	}

	accounts := make(model.Accounts, 0)
	if err := json.Unmarshal(accountsJSON, &accounts); err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	repo.cache[KeyAccountsFindAll].SetDefault(field, accounts)
	return accounts, nil
}

func (repo *redisAccountRepo) Update(account *model.Account) error {
	if err := repo.next.Update(account); err != nil {
		log.Println(err)
		return err
	}

	if err := repo.clearAllFindListCache(); err != nil {
		log.Println(err)
		return err
	}

	field := fmt.Sprintf("%v", account.AccountID)
	if _, err := repo.do("HDEL", KeyAccountsFind, field); err != nil {
		log.Println(err)
		return apperror.InternalServerError
	}

	repo.cache[KeyAccountsFind].Delete(field)
	return nil
}
