package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/s-larionov/process-manager"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/goverland-labs/goverland-core-feed/internal/feedevent"
	"github.com/goverland-labs/goverland-core-feed/internal/pubsub"
	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"

	"github.com/goverland-labs/goverland-core-feed/internal/config"
	"github.com/goverland-labs/goverland-core-feed/internal/item"
	"github.com/goverland-labs/goverland-core-feed/internal/subscriber"
	"github.com/goverland-labs/goverland-core-feed/internal/subscription"
	"github.com/goverland-labs/goverland-core-feed/pkg/grpcsrv"
	"github.com/goverland-labs/goverland-core-feed/pkg/health"
	"github.com/goverland-labs/goverland-core-feed/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App
	db      *gorm.DB

	subscribers      *subscriber.Service
	subscriptions    *subscription.Service
	itemService      *item.Service
	feedEventService *feedevent.Service
}

func NewApplication(cfg config.App) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		sigChan: sigChan,
		cfg:     cfg,
		manager: process.NewManager(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) bootstrap() error {
	initializers := []func() error{
		a.initDB,

		// Init Dependencies
		a.initServices,

		// Init Workers: Application
		// TODO

		// Init Workers: System
		a.initPrometheusWorker,
		a.initHealthWorker,
	}

	for _, initializer := range initializers {
		if err := initializer(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) initDB() error {
	db, err := gorm.Open(postgres.Open(a.cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	ps, err := db.DB()
	if err != nil {
		return err
	}
	ps.SetMaxOpenConns(a.cfg.DB.MaxOpenConnections)

	a.db = db
	if a.cfg.DB.Debug {
		a.db = db.Debug()
	}

	return err
}

func (a *Application) initServices() error {
	nc, err := nats.Connect(
		a.cfg.Nats.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(a.cfg.Nats.MaxReconnects),
		nats.ReconnectWait(a.cfg.Nats.ReconnectTimeout),
	)
	if err != nil {
		return err
	}

	pb, err := natsclient.NewPublisher(nc)
	if err != nil {
		return err
	}

	if err = a.initSubscribers(); err != nil {
		return err
	}
	if err = a.initSubscription(); err != nil {
		return err
	}

	err = a.initDataConsumers(nc, pb)
	if err != nil {
		return fmt.Errorf("init dao: %w", err)
	}

	err = a.initAPI()
	if err != nil {
		return fmt.Errorf("init API: %w", err)
	}

	return nil
}

func (a *Application) initDataConsumers(nc *nats.Conn, pb *natsclient.Publisher) error {
	feedItemsNotifier := pubsub.NewPubSub[string](1000) // TODO: const
	repo := item.NewRepo(a.db)

	service, err := item.NewService(repo, pb, a.subscribers, a.subscriptions, feedItemsNotifier)
	if err != nil {
		return fmt.Errorf("item service: %w", err)
	}

	feedEventService := feedevent.NewService(feedItemsNotifier, service)

	a.itemService = service
	a.feedEventService = feedEventService

	dc, err := item.NewDaoConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("item dao consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("item-dao-consumer", dc.Start))

	pc, err := item.NewProposalConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("item proposal consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("item-proposal-consumer", pc.Start))

	dlc, err := item.NewDelegatesConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("item proposal consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("item-delegate-consumer", dlc.Start))

	return nil
}

// todo: move exclude path to config?
func (a *Application) initAPI() error {
	authInterceptor := grpcsrv.NewAuthInterceptor(a.subscribers)
	srv := grpcsrv.NewGrpcServer(
		[]string{
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
			"/feedpb.Subscriber/Create",
			"/feedpb.Feed/GetByFilter",
		},
		authInterceptor.AuthAndIdentifyTickerFunc,
	)

	feedpb.RegisterSubscriberServer(srv, subscriber.NewServer(a.subscribers))
	feedpb.RegisterSubscriptionServer(srv, subscription.NewServer(a.subscriptions))
	feedpb.RegisterFeedServer(srv, item.NewServer(a.itemService))
	feedpb.RegisterFeedEventsServer(srv, feedevent.NewServer(a.feedEventService))

	a.manager.AddWorker(grpcsrv.NewGrpcServerWorker("API", srv, a.cfg.InternalAPI.Bind))

	return nil
}

func (a *Application) initSubscribers() error {
	repo := subscriber.NewRepo(a.db)
	cache := subscriber.NewCache()
	service, err := subscriber.NewService(repo, cache)
	if err != nil {
		return fmt.Errorf("subsceiber service: %w", err)
	}
	a.subscribers = service

	return nil
}

func (a *Application) initSubscription() error {
	repo := subscription.NewRepo(a.db)
	cache := subscription.NewCache()
	service, err := subscription.NewService(repo, cache)
	if err != nil {
		return fmt.Errorf("subscription service: %w", err)
	}
	a.subscriptions = service

	return nil
}

func (a *Application) initPrometheusWorker() error {
	srv := prometheus.NewServer(a.cfg.Prometheus.Listen, "/metrics")
	a.manager.AddWorker(process.NewServerWorker("prometheus", srv))

	return nil
}

func (a *Application) initHealthWorker() error {
	srv := health.NewHealthCheckServer(a.cfg.Health.Listen, "/status", health.DefaultHandler(a.manager))
	a.manager.AddWorker(process.NewServerWorker("health", srv))

	return nil
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	a.manager.AwaitAll()
}
