package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/atletaid/go-template/src/common/apperror"
	"github.com/atletaid/go-template/src/model"
	"github.com/atletaid/go-template/src/module/restaurant"
	"github.com/lib/pq"
	"github.com/tokopedia/sqlt"
)

type postgreRestaurantRepo struct {
	DbMaster *sqlt.DB
	DbSlave  *sqlt.DB
	Timeout  time.Duration
}

func NewRestaurantRepository(dbMaster *sqlt.DB, dbSlave *sqlt.DB, timeout time.Duration) restaurant.RestaurantRepository {
	return &postgreRestaurantRepo{
		DbMaster: dbMaster,
		DbSlave:  dbSlave,
		Timeout:  timeout,
	}
}

func (repo *postgreRestaurantRepo) CreateRestaurant(restaurant *model.Restaurant) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()
	query := `
		INSERT INTO
			ms_restaurant
		(
			restaurant_name,
			restaurant_time_minute,
			restaurant_price,
			position_lat,
			position_long,
			restaurant_city,
			restaurant_image,
			restaurant_description,
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
			restaurant_id
	`

	var lastInsertID int64
	log.Println(restaurant)
	err := repo.DbMaster.QueryRowContext(
		ctx,
		query,
		restaurant.RestaurantName,
		restaurant.RestaurantTimeMinute,
		restaurant.RestaurantPrice,
		restaurant.PositionLat,
		restaurant.PositionLong,
		restaurant.RestaurantCity,
		restaurant.RestaurantImage,
		restaurant.RestaurantDescription,
	).Scan(&lastInsertID)
	if err != nil {
		log.Println(err)
		return 0, apperror.InternalServerError
	}

	return lastInsertID, nil
}

func (repo *postgreRestaurantRepo) FindRestaurantByID(restaurantID int64) (*model.Restaurant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			restaurant_id,
			restaurant_name,
			restaurant_time_minute,
			restaurant_price,
			position_lat,
			position_long,
			restaurant_city,
			restaurant_image,
			restaurant_description,
			created_at
		FROM
			ms_restaurant
		WHERE
			restaurant_id = $1
	`

	var (
		rtrestaurantID          sql.NullInt64
		rtrestaurantName        sql.NullString
		rtrestaurantTimeMinute  sql.NullInt64
		rtrestaurantPrice       sql.NullInt64
		rtPositionLat           sql.NullFloat64
		rtPositionLong          sql.NullFloat64
		rtrestaurantCity        sql.NullString
		rtrestaurantImage       sql.NullString
		rtrestaurantDescription sql.NullString
		rtCreatedAt             pq.NullTime
	)

	err := repo.DbSlave.QueryRowContext(ctx, query, restaurantID).Scan(
		&rtrestaurantID,
		&rtrestaurantName,
		&rtrestaurantTimeMinute,
		&rtrestaurantPrice,
		&rtPositionLat,
		&rtPositionLong,
		&rtrestaurantCity,
		&rtrestaurantImage,
		&rtrestaurantDescription,
		&rtCreatedAt,
	)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.RestaurantNotExists
	}

	if err != nil {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	restaurant := model.Restaurant{
		RestaurantID:          rtrestaurantID.Int64,
		RestaurantName:        rtrestaurantName.String,
		RestaurantTimeMinute:  int(rtrestaurantTimeMinute.Int64),
		RestaurantPrice:       int(rtrestaurantPrice.Int64),
		PositionLat:           rtPositionLat.Float64,
		PositionLong:          rtPositionLong.Float64,
		RestaurantCity:        rtrestaurantCity.String,
		RestaurantImage:       rtrestaurantImage.String,
		RestaurantDescription: rtrestaurantDescription.String,
		CreatedAt:             rtCreatedAt.Time,
	}

	return &restaurant, nil
}

func (repo *postgreRestaurantRepo) FindAllRestaurants() (model.Restaurants, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
	SELECT
		restaurant_id,
		restaurant_name,
		restaurant_time_minute,
		restaurant_price,
		position_lat,
		position_long,
		restaurant_city,
		restaurant_image,
		restaurant_description,
		created_at
	FROM
		ms_restaurant	
	`

	rows, err := repo.DbSlave.QueryContext(ctx, query)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	restaurants := make(model.Restaurants, 0)
	for rows.Next() {
		var (
			rtrestaurantID          sql.NullInt64
			rtrestaurantName        sql.NullString
			rtrestaurantTimeMinute  sql.NullInt64
			rtrestaurantPrice       sql.NullInt64
			rtPositionLat           sql.NullFloat64
			rtPositionLong          sql.NullFloat64
			rtrestaurantCity        sql.NullString
			rtrestaurantImage       sql.NullString
			rtrestaurantDescription sql.NullString
			rtCreatedAt             pq.NullTime
		)

		if err := rows.Scan(
			&rtrestaurantID,
			&rtrestaurantName,
			&rtrestaurantTimeMinute,
			&rtrestaurantPrice,
			&rtPositionLat,
			&rtPositionLong,
			&rtrestaurantCity,
			&rtrestaurantImage,
			&rtrestaurantDescription,
			&rtCreatedAt,
		); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		restaurant := model.Restaurant{
			RestaurantID:          rtrestaurantID.Int64,
			RestaurantName:        rtrestaurantName.String,
			RestaurantTimeMinute:  int(rtrestaurantTimeMinute.Int64),
			RestaurantPrice:       int(rtrestaurantPrice.Int64),
			PositionLat:           rtPositionLat.Float64,
			PositionLong:          rtPositionLong.Float64,
			RestaurantCity:        rtrestaurantCity.String,
			RestaurantImage:       rtrestaurantImage.String,
			RestaurantDescription: rtrestaurantDescription.String,
			CreatedAt:             rtCreatedAt.Time,
		}

		restaurants = append(restaurants, &restaurant)
	}

	return restaurants, nil
}

func (repo *postgreRestaurantRepo) DeleteRestaurantID(restaurantID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		DELETE FROM
			ms_restaurant
		WHERE
			restaurant_id = $1
	`

	_, err := repo.DbMaster.ExecContext(ctx, query, restaurantID)
	if err != nil {
		log.Println(err)
		return apperror.InternalServerError
	}

	return nil
}

func (repo *postgreRestaurantRepo) FindByLocation(cityName string) (model.Restaurants, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
	SELECT
		restaurant_id,
		restaurant_name,
		restaurant_time_minute,
		restaurant_price,
		position_lat,
		position_long,
		restaurant_city,
		restaurant_image,
		restaurant_description,
		created_at
	FROM
		ms_restaurant	
	WHERE
		restaurant_city = $1
	`

	rows, err := repo.DbSlave.QueryContext(ctx, query, cityName)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, apperror.InternalServerError
	}

	restaurants := make(model.Restaurants, 0)
	for rows.Next() {
		var (
			rtrestaurantID          sql.NullInt64
			rtrestaurantName        sql.NullString
			rtrestaurantTimeMinute  sql.NullInt64
			rtrestaurantPrice       sql.NullInt64
			rtPositionLat           sql.NullFloat64
			rtPositionLong          sql.NullFloat64
			rtrestaurantCity        sql.NullString
			rtrestaurantImage       sql.NullString
			rtrestaurantDescription sql.NullString
			rtCreatedAt             pq.NullTime
		)

		if err := rows.Scan(
			&rtrestaurantID,
			&rtrestaurantName,
			&rtrestaurantTimeMinute,
			&rtrestaurantPrice,
			&rtPositionLat,
			&rtPositionLong,
			&rtrestaurantCity,
			&rtrestaurantImage,
			&rtrestaurantDescription,
			&rtCreatedAt,
		); err != nil {
			log.Println(err)
			return nil, apperror.InternalServerError
		}

		restaurant := model.Restaurant{
			RestaurantID:          rtrestaurantID.Int64,
			RestaurantName:        rtrestaurantName.String,
			RestaurantTimeMinute:  int(rtrestaurantTimeMinute.Int64),
			RestaurantPrice:       int(rtrestaurantPrice.Int64),
			PositionLat:           rtPositionLat.Float64,
			PositionLong:          rtPositionLong.Float64,
			RestaurantCity:        rtrestaurantCity.String,
			RestaurantImage:       rtrestaurantImage.String,
			RestaurantDescription: rtrestaurantDescription.String,
			CreatedAt:             rtCreatedAt.Time,
		}

		restaurants = append(restaurants, &restaurant)
	}

	return restaurants, nil
}
