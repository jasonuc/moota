package events

import (
	"context"
	"database/sql"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-http/v2/pkg/http"
	wmsql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type SeedPlanted struct {
	UserID string `json:"user_id"`
	SeedID string `json:"seed_id"`
}

type SeedGenerated struct {
	UserID string `json:"user_id"`
}

type StatUpdated struct {
	Plant *models.PlantCount `json:"plant,omitempty"`
	Seed  *models.SeedCount  `json:"seed,omitempty"`
}

type Routers struct {
	EventsRouter *message.Router
	SSERouter    http.SSERouter
	EventBus     *cqrs.EventBus
}

func NewPubSub() (message.Publisher, message.Subscriber) {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{
			OutputChannelBuffer: 10000,
		},
		watermill.NewStdLogger(true, true),
	)
	return pubSub, pubSub
}

func NewPostgresPublisher(db *sql.DB) *wmsql.Publisher {
	publisher, err := wmsql.NewPublisher(
		db,
		wmsql.PublisherConfig{
			SchemaAdapter: wmsql.DefaultPostgreSQLSchema{},
		},
		watermill.NewStdLogger(true, true),
	)
	if err != nil {
		panic(err)
	}
	return publisher
}

func NewPostgresSubscriber(db *sql.DB) (message.Subscriber, error) {
	subscriber, err := wmsql.NewSubscriber(
		db,
		wmsql.SubscriberConfig{
			SchemaAdapter:    wmsql.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   wmsql.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: true,
		},
		watermill.NewStdLogger(true, false),
	)
	return subscriber, err
}

func NewEventBus(pubSub message.Publisher) (*cqrs.EventBus, error) {
	return cqrs.NewEventBusWithConfig(
		pubSub,
		cqrs.EventBusConfig{
			GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
				return params.EventName, nil
			},
			Marshaler: cqrs.JSONMarshaler{},
			Logger:    watermill.NewStdLogger(true, true),
		},
	)
}

func NewRouters(store *store.Store, db *sql.DB) (*Routers, error) {

	// publisher := NewPostgresPublisher(db)

	publisher, subscriber := NewPubSub()

	eventBus, err := NewEventBus(publisher)
	if err != nil {
		return nil, err
	}

	eventsRouter, err := NewEventRouter2(subscriber, store, eventBus)
	if err != nil {
		return nil, err
	}

	// subscriber, err := NewPostgresSubscriber(db)
	// if err != nil {
	// 	return nil, err
	// }
	sseRouter, err := NewSseRouter(subscriber)
	if err != nil {
		return nil, err
	}

	return &Routers{
		EventsRouter: eventsRouter,
		SSERouter:    sseRouter,
		EventBus:     eventBus,
	}, nil
}

func NewSseRouter(subscriber message.Subscriber) (http.SSERouter, error) {
	sseRouter, err := http.NewSSERouter(
		http.SSERouterConfig{
			UpstreamSubscriber: subscriber,
		},
		watermill.NewStdLogger(true, true),
	)
	return sseRouter, err
}

func NewEventRouter(db *sql.DB, store *store.Store, eventBus *cqrs.EventBus) (*message.Router, error) {
	logger := watermill.NewStdLogger(true, true)
	eventsRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	eventsRouter.AddMiddleware(middleware.Recoverer)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		eventsRouter,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return NewPostgresSubscriber(db)
			},
			Marshaler: cqrs.JSONMarshaler{},
			Logger:    logger,
		},
	)
	if err != nil {
		return nil, err
	}

	err = eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"UpdatePlantCount",
			func(ctx context.Context, event *SeedPlanted) error {
				count, err := store.Plant.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				count2, err := store.Seed.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				statsUpdated := StatUpdated{
					Plant: count,
					Seed:  count2,
				}

				return eventBus.Publish(ctx, statsUpdated)
			},
		),
		cqrs.NewEventHandler(
			"UpdateSeedCount",
			func(ctx context.Context, event *SeedGenerated) error {
				count, err := store.Seed.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				count2, err := store.Plant.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				statsUpdated := StatUpdated{
					Seed:  count,
					Plant: count2,
				}

				return eventBus.Publish(ctx, statsUpdated)
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return eventsRouter, nil
}
func NewEventRouter2(subscriber message.Subscriber, store *store.Store, eventBus *cqrs.EventBus) (*message.Router, error) {
	logger := watermill.NewStdLogger(true, true)
	eventsRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	eventsRouter.AddMiddleware(middleware.Recoverer)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		eventsRouter,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return subscriber, nil
			},
			Marshaler: cqrs.JSONMarshaler{},
			Logger:    logger,
		},
	)
	if err != nil {
		return nil, err
	}

	err = eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"UpdatePlantCount",
			func(ctx context.Context, event *SeedPlanted) error {
				count, err := store.Plant.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				count2, err := store.Seed.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				statsUpdated := StatUpdated{
					Plant: count,
					Seed:  count2,
				}

				return eventBus.Publish(ctx, statsUpdated)
			},
		),
		cqrs.NewEventHandler(
			"UpdateSeedCount",
			func(ctx context.Context, event *SeedGenerated) error {
				count, err := store.Seed.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				count2, err := store.Plant.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				statsUpdated := StatUpdated{
					Seed:  count,
					Plant: count2,
				}

				return eventBus.Publish(ctx, statsUpdated)
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return eventsRouter, nil
}
