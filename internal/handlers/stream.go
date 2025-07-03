package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jasonuc/moota/internal/events"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type statsStream struct {
	repo *store.Store
}

func (s *statsStream) InitialStreamResponse(w http.ResponseWriter, r *http.Request) (response interface{}, ok bool) {

	plant, err := s.repo.Plant.GetTotalCount(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not get post"))
		return nil, false
	}
	seed, err := s.repo.Seed.GetTotalCount(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not get post"))
		return nil, false
	}
	stats := models.Stats{
		Plant: plant,
		Seed:  seed,
	}

	return stats, true
}

func (s *statsStream) NextStreamResponse(r *http.Request, msg *message.Message) (response interface{}, ok bool) {

	var event events.StatUpdated
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		fmt.Println("cannot unmarshal: " + err.Error())
		return nil, false
	}

	stats := models.Stats{
		Plant: event.Plant,
		Seed:  event.Seed,
	}
	return stats, true
}
