package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/sportivaid/go-template/src/common/apperror"
	"github.com/sportivaid/go-template/src/model"
	"github.com/sportivaid/go-template/src/module/account"
	"github.com/tokopedia/sqlt"
)

type postgreAccountRepo struct {
	DbMaster *sqlt.DB
	DbSlave  *sqlt.DB
	Timeout  time.Duration
}

func NewAccountRepository(dbMaster *sqlt.DB, dbSlave *sqlt.DB, timeout time.Duration) account.AccountRepository {
	return &postgreAccountRepo{
		DbMaster: dbMaster,
		DbSlave:  dbSlave,
		Timeout:  timeout,
	}
}

func (repo *postgreAccountRepo) Create(account *model.Account) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO
			accounts
		(
			user_email,
			user_fullname,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			now(),
			now()
		)
		RETURNING
			account_id
	`

	var lastInsertID int64

	err := repo.DbMaster.QueryRowContext(
		ctx,
		query,
		account.Email,
		account.Fullname,
	).Scan(&lastInsertID)
	if err != nil {
		log.Println(err)
		return 0, apperror.InternalServerError
	}

	return lastInsertID, nil
}

func (repo *postgreAccountRepo) FindByID(accountID int64) (*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			account_id,
			user_email,
			user_fullname,
			created_at,
			updated_at
		FROM
			accounts
		WHERE
			account_id = $1
	`

	var (
		aAccountID sql.NullInt64
		aEmail     sql.NullString
		aFullname  sql.NullString
		aCreatedAt pq.NullTime
		aUpdatedAt pq.NullTime
	)

	err := repo.DbSlave.QueryRowContext(ctx, query, accountID).Scan(
		&aAccountID,
		&aEmail,
		&aFullname,
		&aCreatedAt,
		&aUpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.AccountNotExists
	}

	if err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	account := model.Account{
		AccountID: aAccountID.Int64,
		Email:     aEmail.String,
		Fullname:  aFullname.String,
		CreatedAt: aCreatedAt.Time,
		UpdatedAt: aUpdatedAt.Time,
	}

	return &account, nil
}

func (repo *postgreAccountRepo) FindAll() (model.Accounts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			account_id,
			user_email,
			user_fullname,
			created_at,
			updated_at
		FROM
			accounts
	`

	rows, err := repo.DbSlave.QueryContext(ctx, query)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	accounts := make(model.Accounts, 0)
	for rows.Next() {
		var (
			aAccountID sql.NullInt64
			aEmail     sql.NullString
			aFullname  sql.NullString
			aCreatedAt pq.NullTime
			aUpdatedAt pq.NullTime
		)

		if err := rows.Scan(
			&aAccountID,
			&aEmail,
			&aFullname,
			&aCreatedAt,
			&aUpdatedAt,
		); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		account := model.Account{
			AccountID: aAccountID.Int64,
			Email:     aEmail.String,
			Fullname:  aFullname.String,
			CreatedAt: aCreatedAt.Time,
			UpdatedAt: aUpdatedAt.Time,
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (repo *postgreAccountRepo) Update(account *model.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		UPDATE
			accounts
		SET
			user_email = $2,
			user_fullname = $3,
			updated_at = now()
		WHERE
			account_id = $1
	`

	if _, err := repo.DbMaster.ExecContext(
		ctx,
		query,
		account.AccountID,
		account.Email,
		account.Fullname,
	); err != nil {
		log.Println(err)
		return apperror.InternalServerError
	}

	return nil
}
