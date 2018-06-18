package repository

import (
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	cache "github.com/patrickmn/go-cache"
)

func NewPool(host string, dialTimeout time.Duration, idleTimeout time.Duration, poolSize int) (*redigo.Pool, error) {
	pool := redigo.Pool{
		MaxActive:   poolSize,
		MaxIdle:     poolSize,
		IdleTimeout: idleTimeout,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", host, redigo.DialConnectTimeout(dialTimeout))
			if err != nil {
				log.Println(err)
				return nil, err
			}
			return c, err
		},
	}

	if _, err := pool.Dial(); err != nil {
		pool.Close()
		log.Println(err)
		return nil, err
	}

	return &pool, nil
}

func NewAccountCache(cExpiration time.Duration, cIntervalPurges time.Duration) map[string]*cache.Cache {
	return map[string]*cache.Cache{
		KeyAccountsFindAll: cache.New(cExpiration*time.Minute, cIntervalPurges*time.Minute),
		KeyAccountsFind:    cache.New(cExpiration*time.Minute, cIntervalPurges*time.Minute),
	}
}
