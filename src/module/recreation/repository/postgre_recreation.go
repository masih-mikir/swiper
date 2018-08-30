package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/atletaid/go-template/src/common/apperror"
	"github.com/atletaid/go-template/src/model"
	"github.com/atletaid/go-template/src/module/recreation"
	"github.com/lib/pq"
	"github.com/tokopedia/sqlt"
)

type postgreRecreationRepo struct {
	DbMaster *sqlt.DB
	DbSlave  *sqlt.DB
	Timeout  time.Duration
}

func NewRecreationRepository(dbMaster *sqlt.DB, dbSlave *sqlt.DB, timeout time.Duration) recreation.RecreationRepository {
	return &postgreRecreationRepo{
		DbMaster: dbMaster,
		DbSlave:  dbSlave,
		Timeout:  timeout,
	}
}

func (repo *postgreRecreationRepo) CreateRecreation(recreation *model.Recreation) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO
			ms_recreation
		(
			recreation_name,
			recreation_time_minute,
			recreation_price,
			position_lat,
			position_long,
			recreation_city,
			recreation_image,
			recreation_description,
			created_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			now()
		)
		RETURNING
			recreation_id
	`

	var lastInsertID int64

	err := repo.DbMaster.QueryRowContext(
		ctx,
		query,
		recreation.RecreationName,
		recreation.RecreationTimeMinute,
		recreation.RecreationPrice,
		recreation.PositionLat,
		recreation.PositionLong,
		recreation.RecreationCity,
		recreation.RecreationImage,
		recreation.RecreationDescription,
	).Scan(&lastInsertID)
	if err != nil {
		log.Println(err)
		return 0, apperror.InternalServerError
	}

	return lastInsertID, nil
}

func (repo *postgreRecreationRepo) FindRecreationByID(recreationID int64) (*model.Recreation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			recreation_id,
			recreation_name,
			recreation_time_minute,
			recreation_price,
			position_lat,
			position_long,
			recreation_city,
			recreation_image,
			recreation_description,
			created_at
		FROM
			ms_recreation
		WHERE
			recreation_id = $1
	`

	var (
		rRecreationID          sql.NullInt64
		rRecreationName        sql.NullString
		rRecreationTimeMinute  sql.NullInt64
		rRecreationPrice       sql.NullInt64
		rPositionLat           sql.NullFloat64
		rPositionLong          sql.NullFloat64
		rRecreationCity        sql.NullString
		rRecreationImage       sql.NullString
		rRecreationDescription sql.NullString
		rCreatedAt             pq.NullTime
	)

	err := repo.DbSlave.QueryRowContext(ctx, query, recreationID).Scan(
		&rRecreationID,
		&rRecreationName,
		&rRecreationTimeMinute,
		&rRecreationPrice,
		&rPositionLat,
		&rPositionLong,
		&rRecreationCity,
		&rRecreationImage,
		&rRecreationDescription,
		&rCreatedAt,
	)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.RecreationNotExists
	}

	if err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	recreation := model.Recreation{
		RecreationID:          rRecreationID.Int64,
		RecreationName:        rRecreationName.String,
		RecreationTimeMinute:  int(rRecreationTimeMinute.Int64),
		RecreationPrice:       int(rRecreationPrice.Int64),
		PositionLat:           rPositionLat.Float64,
		PositionLong:          rPositionLong.Float64,
		RecreationCity:        rRecreationCity.String,
		RecreationDescription: rRecreationDescription.String,
		CreatedAt:             rCreatedAt.Time,
	}

	return &recreation, nil
}

func (repo *postgreRecreationRepo) FindAllRecreations() (model.Recreations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
	SELECT
		recreation_id,
		recreation_name,
		recreation_time_minute,
		recreation_price,
		position_lat,
		position_long,
		recreation_city,
		recreation_image,
		recreation_description,
		created_at
	FROM
		ms_recreation	
	`

	rows, err := repo.DbSlave.QueryContext(ctx, query)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	recreations := make(model.Recreations, 0)
	for rows.Next() {
		var (
			rRecreationID          sql.NullInt64
			rRecreationName        sql.NullString
			rRecreationTimeMinute  sql.NullInt64
			rRecreationPrice       sql.NullInt64
			rPositionLat           sql.NullFloat64
			rPositionLong          sql.NullFloat64
			rRecreationCity        sql.NullString
			rRecreationImage       sql.NullString
			rRecreationDescription sql.NullString
			rCreatedAt             pq.NullTime
		)

		if err := rows.Scan(
			&rRecreationID,
			&rRecreationName,
			&rRecreationTimeMinute,
			&rRecreationPrice,
			&rPositionLat,
			&rPositionLong,
			&rRecreationCity,
			&rRecreationImage,
			&rRecreationDescription,
			&rCreatedAt,
		); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		recreation := model.Recreation{
			RecreationID:          rRecreationID.Int64,
			RecreationName:        rRecreationName.String,
			RecreationTimeMinute:  int(rRecreationTimeMinute.Int64),
			RecreationPrice:       int(rRecreationPrice.Int64),
			PositionLat:           rPositionLat.Float64,
			PositionLong:          rPositionLong.Float64,
			RecreationCity:        rRecreationCity.String,
			RecreationDescription: rRecreationDescription.String,
			CreatedAt:             rCreatedAt.Time,
		}

		recreations = append(recreations, &recreation)
	}

	return recreations, nil
}

func (repo *postgreRecreationRepo) DeleteRecreation(recreationID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		DELETE FROM
			ms_recreation
		WHERE
			recreation_id = $1
	`

	_, err := repo.DbMaster.ExecContext(ctx, query, recreationID)
	if err != nil {
		log.Println(err)
		return apperror.InternalServerError
	}

	return nil
}

func (repo *postgreRecreationRepo) FindByLocation(cityName string) (model.Recreations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
	SELECT
		recreation_id,
		recreation_name,
		recreation_time_minute,
		recreation_price,
		position_lat,
		position_long,
		recreation_city,
		recreation_image,
		recreation_description,
		created_at
	FROM
		ms_recreation	
	WHERE
		recreation_city = $1
	`

	rows, err := repo.DbSlave.QueryContext(ctx, query, cityName)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	recreations := make(model.Recreations, 0)
	for rows.Next() {
		var (
			rRecreationID          sql.NullInt64
			rRecreationName        sql.NullString
			rRecreationTimeMinute  sql.NullInt64
			rRecreationPrice       sql.NullInt64
			rPositionLat           sql.NullFloat64
			rPositionLong          sql.NullFloat64
			rRecreationCity        sql.NullString
			rRecreationImage       sql.NullString
			rRecreationDescription sql.NullString
			rCreatedAt             pq.NullTime
		)

		if err := rows.Scan(
			&rRecreationID,
			&rRecreationName,
			&rRecreationTimeMinute,
			&rRecreationPrice,
			&rPositionLat,
			&rPositionLong,
			&rRecreationCity,
			&rRecreationImage,
			&rRecreationDescription,
			&rCreatedAt,
		); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		recreation := model.Recreation{
			RecreationID:          rRecreationID.Int64,
			RecreationName:        rRecreationName.String,
			RecreationTimeMinute:  int(rRecreationTimeMinute.Int64),
			RecreationPrice:       int(rRecreationPrice.Int64),
			PositionLat:           rPositionLat.Float64,
			PositionLong:          rPositionLong.Float64,
			RecreationCity:        rRecreationCity.String,
			RecreationDescription: rRecreationDescription.String,
			CreatedAt:             rCreatedAt.Time,
		}

		recreations = append(recreations, &recreation)
	}

	return recreations, nil
}
