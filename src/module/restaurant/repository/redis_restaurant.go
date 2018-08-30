package repository

import (
	"log"

	"github.com/atletaid/go-template/src/common/apperror"
	"github.com/atletaid/go-template/src/model"
	"github.com/atletaid/go-template/src/module/restaurant"
	redigo "github.com/gomodule/redigo/redis"
	cache "github.com/patrickmn/go-cache"
)

const (
	KeyRestaurantsFindAll = "restaurants:find_all"
	KeyRestaurantsFind    = "restaurants:find"
)

type redisRestaurantRepo struct {
	cache map[string]*cache.Cache
	pool  *redigo.Pool
	next  restaurant.RestaurantRepository
}

func NewMiddlewareRestaurantRepository(cache map[string]*cache.Cache, pool *redigo.Pool, next restaurant.RestaurantRepository) restaurant.RestaurantRepository {
	return &redisRestaurantRepo{
		cache: cache,
		pool:  pool,
		next:  next,
	}
}

func (repo *redisRestaurantRepo) do(command string, args ...interface{}) (reply interface{}, err error) {
	conn := repo.pool.Get()
	defer conn.Close()

	return conn.Do(command, args...)
}

func (repo *redisRestaurantRepo) clearAllFindListCache() error {
	keys := []interface{}{
		KeyRestaurantsFindAll,
	}

	if _, err := repo.do("DEL", keys...); err != nil {
		log.Println(err)
		return apperror.InternalServerError
	}

	repo.cache[KeyRestaurantsFindAll].Flush()
	return nil
}

func (repo *redisRestaurantRepo) CreateRestaurant(restaurant *model.Restaurant) (int64, error) {
	return repo.next.CreateRestaurant(restaurant)
}

func (repo *redisRestaurantRepo) FindRestaurantByID(restaurantID int64) (*model.Restaurant, error) {
	return repo.next.FindRestaurantByID(restaurantID)
}

func (repo *redisRestaurantRepo) FindAllRestaurants() (model.Restaurants, error) {
	return repo.next.FindAllRestaurants()
}

func (repo *redisRestaurantRepo) FindByLocation(cityName string) (model.Restaurants, error) {
	return repo.next.FindByLocation(cityName)
}

func (repo *redisRestaurantRepo) DeleteRestaurantID(restaurantID int64) error {
	return repo.next.DeleteRestaurantID(restaurantID)
}
