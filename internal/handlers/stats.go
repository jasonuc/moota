package handlers

import (
	"net/http"

	"github.com/jasonuc/moota/internal/events"
	"github.com/jasonuc/moota/internal/store"
	"github.com/jasonuc/moota/internal/websocket"
)

type StatHandler struct {
	SseHandler http.HandlerFunc
}

func NewWssStatHandler(broadcaster websocket.Broadcaster, store *store.Store) *StatHandler {
	statsHandler := events.NewSockServer(broadcaster, store)
	return &StatHandler{
		SseHandler: statsHandler,
	}
}
