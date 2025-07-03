package handlers

import (
	"net/http"

	watermillhttp "github.com/ThreeDotsLabs/watermill-http/v2/pkg/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/jasonuc/moota/internal/events"
	"github.com/jasonuc/moota/internal/store"
)

type StatHandler struct {
	SseHandler http.HandlerFunc
}

func NewStatHandler(routers watermillhttp.SSERouter, store *store.Store) *StatHandler {
	marshaler := cqrs.JSONMarshaler{}
	topic := marshaler.Name(events.StatUpdated{})
	statsHandler := routers.AddHandler(topic, &statsStream{repo: store})
	return &StatHandler{
		SseHandler: statsHandler,
	}
}
