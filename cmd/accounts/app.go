package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/sportivaid/go-template/config"
	"github.com/sportivaid/go-template/src/module/account"
	"github.com/sportivaid/go-template/src/module/account/delivery"
	"github.com/sportivaid/go-template/src/module/account/repository"
	"github.com/tokopedia/sqlt"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ldate)

	//Init Config
	cfg, ok := config.InitConfig([]string{"files/etc/config"}...)
	if !ok {
		fmt.Println("Error opening config files")
		return
	}

	flag.Parse()

	// Init PostgreSQL Database
	dbMaster, err := sqlt.Open("postgres", cfg.Account.MasterDB)
	if err != nil {
		log.Println("Error opening database : ", err)
		return
	}

	// Init Inmemory & Redis Cache
	accountCache := repository.NewAccountCache(cfg.InMemory.DefaultExpiration, cfg.InMemory.IntervalPurges)
	redisPool, err := repository.NewPool(cfg.Redis.Host, cfg.Redis.DialTimeout*time.Second, cfg.Redis.IdleTimeout*time.Second, cfg.Redis.PoolSize)
	if err != nil {
		log.Println(err)
		return
	}

	accountRepo := repository.NewAccountRepository(dbMaster, dbMaster, cfg.Server.DBTimeout*time.Second)
	accountRepo = repository.NewMiddlewareAccountRepository(accountCache, redisPool, accountRepo)
	accountUsecase := account.NewAccountUsecase(accountRepo)
	router := delivery.NewAccountHandler(accountUsecase)
	router.Run(cfg.Account.Port)
}
