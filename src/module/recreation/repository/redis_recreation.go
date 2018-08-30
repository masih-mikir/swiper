package repository

import (
	"log"

	"github.com/atletaid/go-template/src/common/apperror"
	"github.com/atletaid/go-template/src/model"
	"github.com/atletaid/go-template/src/module/recreation"
	redigo "github.com/gomodule/redigo/redis"
	cache "github.com/patrickmn/go-cache"
)

const (
	KeyRecreationsFindAll = "recreations:find_all"
	KeyRecreationsFind    = "recreations:find"
)

type redisRecreationRepo struct {
	cache map[string]*cache.Cache
	pool  *redigo.Pool
	next  recreation.RecreationRepository
}

func NewMiddlewareRecreationRepository(cache map[string]*cache.Cache, pool *redigo.Pool, next recreation.RecreationRepository) recreation.RecreationRepository {
	return &redisRecreationRepo{
		cache: cache,
		pool:  pool,
		next:  next,
	}
}

func (repo *redisRecreationRepo) do(command string, args ...interface{}) (reply interface{}, err error) {
	conn := repo.pool.Get()
	defer conn.Close()

	return conn.Do(command, args...)
}

func (repo *redisRecreationRepo) clearAllFindListCache() error {
	keys := []interface{}{
		KeyRecreationsFindAll,
	}

	if _, err := repo.do("DEL", keys...); err != nil {
		log.Println(err)
		return apperror.InternalServerError
	}

	repo.cache[KeyRecreationsFindAll].Flush()
	return nil
}

func (repo *redisRecreationRepo) CreateRecreation(recreation *model.Recreation) (int64, error) {
	return repo.next.CreateRecreation(recreation)
}

func (repo *redisRecreationRepo) FindRecreationByID(recreationID int64) (*model.Recreation, error) {
	return repo.next.FindRecreationByID(recreationID)
}

func (repo *redisRecreationRepo) FindAllRecreations() (model.Recreations, error) {
	return repo.next.FindAllRecreations()
}

func (repo *redisRecreationRepo) FindByLocation(cityName string) (model.Recreations, error) {
	return repo.next.FindByLocation(cityName)
}

func (repo *redisRecreationRepo) DeleteRecreation(recreationID int64) error {
	return repo.next.DeleteRecreation(recreationID)
}
