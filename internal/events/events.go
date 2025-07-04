package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	gw "github.com/gorilla/websocket"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
	"github.com/jasonuc/moota/internal/websocket"
)

type SeedPlanted struct {
	UserID string `json:"user_id"`
	SeedID string `json:"seed_id"`
}

type SeedGenerated struct {
	UserID string `json:"user_id"`
}
type StatUpdated struct {
}

type StatUpdatedPayload struct {
	Plant *models.PlantCount `json:"plant,omitempty"`
	Seed  *models.SeedCount  `json:"seed,omitempty"`
}

type Routers struct {
	EventsRouter *message.Router
	EventBus     *cqrs.EventBus
}

func NewPubSub() (message.Publisher, message.Subscriber) {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{
			OutputChannelBuffer: 10000,
		},
		watermill.NopLogger{},
	)
	return pubSub, pubSub
}

func NewEventBus(pubSub message.Publisher) (*cqrs.EventBus, error) {
	return cqrs.NewEventBusWithConfig(
		pubSub,
		cqrs.EventBusConfig{
			GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
				return params.EventName, nil
			},

			Marshaler: cqrs.JSONMarshaler{},
			// Logger:    // watermill.NewStdLogger(true, true),
		},
	)
}

func NewSocketRouters(broadcast websocket.Broadcaster, store *store.Store, db *sql.DB) (*Routers, error) {

	publisher, subscriber := NewPubSub()

	eventBus, err := NewEventBus(publisher)
	if err != nil {
		return nil, err
	}

	eventsRouter, err := NewBroadcastEventRouter(broadcast, subscriber, store, eventBus)
	if err != nil {
		return nil, err
	}

	return &Routers{
		EventsRouter: eventsRouter,
		EventBus:     eventBus,
	}, nil
}

func NewSockServer(manager websocket.Manager, store *store.Store) http.HandlerFunc {
	upgrader := gw.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	h := websocket.ServeWS(
		upgrader,
		websocket.DefaultSetupConn,
		websocket.NewClient,
		func(ctx context.Context, cf context.CancelFunc, c websocket.Client) {
			manager.RegisterClient(ctx, cf, c)
			count, err := store.Plant.GetTotalCount(ctx)
			if err != nil {
				slog.Error("error getting plant count", slog.Any("error", err))
				return
			}
			count2, err := store.Seed.GetTotalCount(ctx)
			if err != nil {
				slog.Error("error getting seed count", slog.Any("error", err))
				return
			}
			statsUpdated := StatUpdatedPayload{
				Plant: count,
				Seed:  count2,
			}
			msg, err := json.Marshal(statsUpdated)
			if err != nil {
				slog.Error("error marshaling stats", slog.Any("error", err))
				return
			}
			if err := c.Conn().WriteMessage(gw.TextMessage, msg); err != nil {
				slog.Error("error writing message", "error", err)
				c.Log(int(slog.LevelError), fmt.Sprintf("error writing message: %v", err))
				return
			}
		},
		func(c websocket.Client) {
			manager.UnregisterClient(c)

		},
		50*time.Second,
		[]websocket.MessageHandler{func(c websocket.Client, b []byte) { c.Write(b) }},
	)

	return h
}
func NewBroadcastEventRouter(broadcaster websocket.Broadcaster, subscriber message.Subscriber, store *store.Store, eventBus *cqrs.EventBus) (*message.Router, error) {
	eventsRouter, err := message.NewRouter(message.RouterConfig{}, watermill.NopLogger{})
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
			// Logger:    logger,
		},
	)
	if err != nil {
		return nil, err
	}

	err = eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"StatUpdated",
			func(ctx context.Context, event *StatUpdated) error {
				slog.Info("StatUpdated")
				count, err := store.Plant.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				count2, err := store.Seed.GetTotalCount(ctx)
				if err != nil {
					return err
				}
				statsUpdated := StatUpdatedPayload{
					Plant: count,
					Seed:  count2,
				}
				msg, _ := json.Marshal(statsUpdated)
				return broadcaster.Broadcast(msg)
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return eventsRouter, nil
}
