package bootstrap

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
	"wallet/internal/config"
	"wallet/internal/drivers/async"
	cch "wallet/internal/drivers/redis"
	"wallet/internal/gateways/amqp"
	"wallet/internal/gateways/rest"
	"wallet/internal/repository"
	"wallet/internal/usecases/exists"
	"wallet/internal/usecases/replenishment"
	"wallet/internal/worker/transaction"
	"wallet/pkg/rabbit"
	"wallet/pkg/tracer"
)

type App struct {
	config   config.Config
	teardown []func()

	rest *rest.Server
	amqp *amqp.AMQP
}

type usecases struct {
	walletExists *exists.UseCase
	transfer     *replenishment.UseCase
}

func (a App) Run(ctx context.Context) {
	go a.rest.Run()

	go func() {
		err := a.amqp.Run()
		if err != nil {
			return
		}
	}()

	<-ctx.Done()

	for _, t := range a.teardown {
		t()
	}
}

func NewApp(cfg config.Config) *App {
	teardown := make([]func(), 0)

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer log.Sync() // nolint:errcheck

	conn := rabbit.NewConnection(cfg)

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.RedisURL(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	redisClient := cch.New(rdb)

	teardown = append(teardown, func() {
		if err := rdb.Close(); err != nil {
			panic(err)
		}
	})

	db, err := initDB(cfg)
	if err != nil {
		panic(err)
	}

	store := repository.NewStore(db)

	teardown = append(teardown, func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	})

	newTracer, err := tracer.NewTracer(cfg.AppLabel, cfg.Jaeger.Host, cfg.Jaeger.Port)
	if err != nil {
		return nil
	}

	uc := buildUsecases(log, cfg, redisClient, store, newTracer, conn)

	restServer := rest.New(cfg, uc.walletExists)

	teardown = append(teardown, func() {
		restServer.Shutdown(context.Background())
	})

	amqpConn := rabbit.NewConnection(cfg)

	teardown = append(teardown, func() {
		if err = amqpConn.Close(); err != nil {
			panic(err)
		}
	})

	tranWorker := transaction.New(log, newTracer, redisClient, store.TransactionRepo(), store.WalletRepo(), store.ProductRepo())

	amqpServer := amqp.NewAMQP(log, cfg, amqpConn, tranWorker)

	return &App{
		config:   cfg,
		teardown: teardown,

		rest: restServer,
		amqp: amqpServer,
	}
}

func buildUsecases(
	log *zap.Logger,
	cfg config.Config,
	cch *cch.Client,
	store *repository.Store,
	tracer *tracer.Tracer,
	conn *rabbit.Connection) *usecases {

	return &usecases{
		walletExists: exists.New(log, cch, store.WalletRepo()),
		transfer: replenishment.New(log, cch, store.WalletRepo(),
			store.TransactionRepo(),
			tracer,
			async.New(log, *cfg.TransactionPub, conn)),
	}
}
